package taobao_client

import (
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/secr3t/taobao-client/model"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	atpClient "github.com/secr3t/taobao-client/atp/client"
	ocClient "github.com/secr3t/taobao-client/openchina/client"
	otClient "github.com/secr3t/taobao-client/ot/client"
	rakutenClient "github.com/secr3t/taobao-client/rakuten/client"
)

var Cache = ttlcache.NewCache()

func init() {
	Cache.SetTTL(10 * time.Minute)
}

type TaobaoClient struct {
	searchClient    *rakutenClient.SearchClient
	otSearchClient  *otClient.SearchClient
	otDetailClient  *otClient.DetailClient
	ocDetailClient  *ocClient.DetailClient
	atpDetailClient *atpClient.DetailClient
}

func NewTaobaoClient(rakutenKey, otKey, atpKey, ocKey string) *TaobaoClient {
	return &TaobaoClient{
		searchClient:    rakutenClient.NewSearchClient(rakutenKey),
		otDetailClient:  otClient.NewDetailClient(otKey),
		otSearchClient:  otClient.NewSearchClient(otKey),
		ocDetailClient:  ocClient.NewDetailClient(ocKey),
		atpDetailClient: atpClient.NewDetailClient(atpKey),
	}
}

func (c *TaobaoClient) Search(uri string) []model.Item {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if value, exists := Cache.Get(uri); exists == nil {
		return value.([]model.Item)
	}

	items := make([]model.Item, 0)
	reqUri, e := url.ParseRequestURI(uri)

	if e != nil {
		return items
	}

	param := otClient.SearchParamFromUri(0, uri)

	if strings.TrimSpace(*param.XmlParameters.ItemTitle) == "" {
		return items
	}

	//if items = ItemsFromAtp(uri); len(items) == 0 {
	if reqUri.Query().Get("ppath") == "" {
		if items = c.searchClient.GetItems(uri); len(items) == 0 {
			items = c.otSearchClient.GetItems(uri)
		}
	} else {
		if items = c.otSearchClient.GetItems(uri); len(items) == 0 {
			items = c.searchClient.GetItems(uri)
		}
	}
	//}

	Cache.Set(uri, items)

	return items
}

func (c *TaobaoClient) DetailChain(items []model.Item) chan model.DetailItem {
	return c.DetailChainWithIds(ItemsToIds(items))
}

func (c *TaobaoClient) DetailChainWithIds(ids []string) chan model.DetailItem {
	var wg sync.WaitGroup
	detailChan := make(chan model.DetailItem, len(ids))
	defer func() {
		wg.Wait()
		close(detailChan)
	}()

	wg.Add(len(ids))

	for _, id := range ids {
		id := id
		go func() {
			var detail *model.DetailItem
			if detail = c.ocDetailClient.GetDetail(id); detail == nil {
				detail = c.atpDetailClient.GetDetail(id)
			}

			if detail != nil {
				detailChan <- *detail
			}
			wg.Done()
		}()
	}

	return detailChan
}

func ItemsToIds(items []model.Item) []string {
	ids := make([]string, 0)

	for _, item := range items {
		ids = append(ids, item.Id)
	}

	return ids
}
