package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/secr3t/taobao-client/atp/model"
	"github.com/secr3t/taobao-client/helper"
	model2 "github.com/secr3t/taobao-client/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const detailApiName = "item_get"

var s = semaphore.NewWeighted(10)

type DetailClient struct {
	apiKey string
}

func NewDetailClient(apiKey string) *DetailClient {
	return &DetailClient{
		apiKey: apiKey,
	}
}

func (c *DetailClient) GetItems(itemIds []string) []model.DetailItem {
	items := make([]model.DetailItem, 0)
	itemChans := make(chan *model.DetailItem)

	wg := sync.WaitGroup{}
	for _, itemId := range itemIds {
		itemId := itemId
		go func() {
			wg.Add(1)
			result, err := c.GetItem(itemId)
			if err == nil {
				result.DetailItem.SetOptions()
				itemChans <- result.DetailItem
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(itemChans)
	}()

	for item := range itemChans {
		items = append(items, *item)
	}

	return items
}

func (c *DetailClient) GetDetail(itemId string) *model2.DetailItem {
	result, err := c.getItem(itemId, true)
	if err != nil {
		return nil
	}
	result.DetailItem.SetOptions()

	return &model2.DetailItem{
		Id: result.DetailItem.NumIid,
		Title: result.DetailItem.Title,
		Price: helper.PriceAsFloat(result.DetailItem.Price),
		ProductURL: result.DetailItem.DetailURL,
		MainImgURL: result.DetailItem.GetMainImg(),
		Images: result.DetailItem.GetItemImgs(),
		DescImages: result.DetailItem.GetDetailImgs(),
		Options: result.DetailItem.Options,
	}
}

func (c *DetailClient) GetDetails(itemIds []string) chan *model.DetailItem {
	itemLen := len(itemIds)

	var wg sync.WaitGroup
	wg.Add(itemLen)

	itemChans := make(chan *model.DetailItem, itemLen)

	delta := time.Millisecond.Milliseconds() * 200

	for idx, itemId := range itemIds {
		itemId := itemId
		duration := int64(idx) * delta
		go func() {
			time.Sleep(time.Duration(duration))
			result, err := c.getItem(itemId, true)
			if err == nil {
				result.DetailItem.SetOptions()
				itemChans <- result.DetailItem
			}
			wg.Done()
		}()
	}

	defer func() {
		go func() {
			wg.Wait()
			close(itemChans)
		}()
	}()

	return itemChans
}

func (c *DetailClient) getItem(itemId string, noCache bool) (model.DetailResult, error) {
	ctx := context.TODO()
	s.Acquire(ctx, 1)
	defer s.Release(1)

	query := c.getDetailQueryParam(itemId, noCache)

	reqUri := GetUri(query)

	res, err := http.Get(reqUri)

	if err != nil {
		return model.DetailResult{}, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var result model.DetailResult

	err = json.Unmarshal(body, &result)

	if err != nil {
		return model.DetailResult{}, err
	}

	if result.IsError() {
		err = errors.New(fmt.Sprintf("%s fetch failed.", itemId))
		return model.DetailResult{}, err
	}

	return result, nil
}

func (c *DetailClient) GetItem(itemId string) (model.DetailResult, error) {
	result, err := c.getItem(itemId, false)

	if err != nil {
		return c.getItem(itemId, true)
	}

	return result, err
}

func (c *DetailClient) getDetailQueryParam(itemId string, noCache bool) string {
	p := url.Values{}

	p.Add("api_name", detailApiName)
	p.Add("route", route)
	p.Add("lang", lang)
	p.Add("is_promotion", "1")
	p.Add("key", c.apiKey)
	p.Add("num_iid", itemId)
	if noCache {
		p.Add("cache", "no")
	}

	return p.Encode()
}
