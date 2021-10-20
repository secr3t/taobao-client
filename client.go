package taobao_client

import (
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/secr3t/taobao-client/model"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	otClient "github.com/secr3t/taobao-client/ot/client"
	rakutenClient "github.com/secr3t/taobao-client/rakuten/client"
)

var Cache = ttlcache.NewCache()

func init() {
	Cache.SetTTL(10 * time.Minute)
}

type TaobaoClient struct {
	idx            int
	mutex          *sync.Mutex
	searchClient   *rakutenClient.SearchClient
	detailClients  []*rakutenClient.DetailClient
	otSearchClient *otClient.SearchClient
}

type RakutenConfig struct {
	Weight int64
	Keys   []string
}

func NewTaobaoClient(otKey string, rakutenConfigs []RakutenConfig, hook func(key string)) *TaobaoClient {
	var allKeys []string
	var detailClients []*rakutenClient.DetailClient
	for _, config := range rakutenConfigs {
		detailClients = append(detailClients, rakutenClient.NewDetailClient(config.Weight, config.Keys, hook).AddOtClient(otKey))
		for _, key := range config.Keys {
			allKeys = append(allKeys, key)
		}
	}
	rakutenClient.InitApiKeys(allKeys...)
	return &TaobaoClient{
		mutex:          &sync.Mutex{},
		searchClient:   rakutenClient.NewSearchClient(),
		otSearchClient: otClient.NewSearchClient(otKey),
		detailClients:  detailClients,
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

	if reqUri.Query().Get("ppath") == "" {
		if items = c.searchClient.GetItems(uri); len(items) == 0 {
			items = c.otSearchClient.GetItems(uri)
		}
	} else {
		if items = c.otSearchClient.GetItems(uri); len(items) == 0 {
			items = c.searchClient.GetItems(uri)
		}
	}

	Cache.Set(uri, items)

	return items
}

func (c *TaobaoClient) DetailChain(items []model.Item) chan model.DetailItem {
	return c.DetailChainWithIds(ItemsToIds(items))
}

func (c *TaobaoClient) GetClient() *rakutenClient.DetailClient {
	c.mutex.Lock()
	defer func() {
		if c.idx+1 == len(c.detailClients) {
			c.idx = 0
		} else {
			c.idx += 1
		}
		c.mutex.Unlock()
	}()

	return c.detailClients[c.idx]
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
			detail, _ = c.GetClient().GetDetail(id)

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
