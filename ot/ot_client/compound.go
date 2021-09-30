package client

import (
	model "github.com/secr3t/taobao-client/ot/ot_model"
	"log"
	"sync"
	"time"
)

const defaultLimit = 200

type CompoundClient struct {
	ApiKey      string
	SearchLimit int
}

func NewCompoundClient(apiKey string, searchLimit int) *CompoundClient {
	if searchLimit > defaultLimit {
		searchLimit = defaultLimit
	}
	return &CompoundClient{
		ApiKey:      apiKey,
		SearchLimit: searchLimit,
	}
}

func (c *CompoundClient) SearchAndGetDetail(param *SearchParam) ([]model.DetailItem, error) {
	limit := c.SearchLimit
	sc := NewSearchClient(c.ApiKey)
	dc := NewDetailClient(c.ApiKey)
	detailItems := make([]model.DetailItem, 0)

	result, err := sc.SearchItems(*param)

	if err != nil {
		return nil, err
	}
	if limit > result.TotalCount {
		limit = result.TotalCount
	}

	limit -= len(result.Items)

	for _, item := range result.Items {
		detailItem, err := dc.GetDetail(item)
		if err != nil {
			log.Printf("id %s fetch failed", detailItem.NumIid)
		} else {
			detailItems = append(detailItems, *detailItem)
		}
	}

	param.Page += result.FrameSize

	for ;limit > 0; limit -= len(result.Items) {
		result, err = sc.SearchItems(*param)
		for _, item := range result.Items {
			detailItem, err := dc.GetDetail(item)
			if err != nil {
				log.Printf("id %s fetch failed", detailItem.NumIid)
			} else {
				detailItems = append(detailItems, *detailItem)
			}
		}

		param.Page += result.FrameSize
	}

	return detailItems, nil
}

func (c *CompoundClient) SearchAndGetDetailsMultiRequestOneTime(param *SearchParam) (chan model.DetailItem, error) {
	items := NewSearchClient(c.ApiKey).SearchTilLimit(param, c.SearchLimit)

	return c.GetDetails(items)
}

func (c *CompoundClient) GetDetails(items []model.Item) (chan model.DetailItem, error) {
	var wg sync.WaitGroup
	itemLen := len(items)
	wg.Add(itemLen)

	detailChan := make(chan model.DetailItem, itemLen)

	for i, item := range items {
		go c.backgroundDetailRequestItem(&wg, item, detailChan, int64(i))
	}

	defer func() {
		go func() {
			wg.Wait()
			close(detailChan)
		}()
	}()

	return detailChan, nil
}

func (c *CompoundClient) backgroundDetailRequestItem(wg *sync.WaitGroup, item model.Item, ch chan model.DetailItem, sleepDelta int64) {
	time.Sleep(time.Millisecond * time.Duration(sleepDelta * 50))
	dc := NewDetailClient(c.ApiKey)

	detail, err := dc.GetDetail(item)

	if err == nil {
		ch <- *detail
	} else {
		log.Println(err)
	}
	wg.Done()
}
