package client

import (
	"context"
	"encoding/json"
	"github.com/secr3t/taobao-client/helper"
	model2 "github.com/secr3t/taobao-client/model"
	otClient "github.com/secr3t/taobao-client/ot/client"
	"github.com/secr3t/taobao-client/rakuten/model"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
)

var (
	s = semaphore.NewWeighted(50)
)

type DetailClient struct {
	ApiKeys  []string
	idx      int64
	OtClient *otClient.DetailClient
}

func NewDetailClient(apiKeys ...string) *DetailClient {
	return &DetailClient{
		ApiKeys: apiKeys,
	}
}

func (c *DetailClient) AddOtClient(otApiKey string) *DetailClient {
	c.OtClient = otClient.NewDetailClient(otApiKey)
	return c
}

func (c *DetailClient) GetRequest(id, api string) *http.Request {
	queryV := url.Values{}
	queryV.Add("num_iid", id)
	uri := GetUri(taobaoApiHost, api, queryV.Encode())

	req, _ := http.NewRequest("GET", uri, nil)

	currentIdx := atomic.LoadInt64(&c.idx)
	apiKey := c.ApiKeys[currentIdx]
	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", taobaoApiHost)

	atomic.AddInt64(&c.idx, 1)
	if c.idx == int64(len(c.ApiKeys)) {
		atomic.StoreInt64(&c.idx, 0)
	}

	return req
}

func (c *DetailClient) GetDetail(id string) *model2.DetailItem {
	var detailItem *model2.DetailItem
	ds := c.getDetail(id)

	if !ds.IsSuccess() {
		if c.OtClient == nil {
			return nil
		}
		if detailItem = c.OtClient.GetDetailBase(id); detailItem == nil {
			return nil
		}
	}

	desc := c.getDesc(id)

	if !desc.IsSuccess() {
		return nil
	}

	sku := c.getSku(id)

	if !sku.IsSuccess() {
		return nil
	}

	if detailItem != nil {
		options := sku.GetOptions()
		var price float64
		for _, option := range options {
			if price == 0 {
				price = option.Price
			} else {
				price = math.Min(price, option.Price)
			}
		}
		detailItem.Options = options
		detailItem.Price = price
		detailItem.DescImages = desc.GetImages()

		return detailItem
	}

	return &model2.DetailItem{
		Id:         strconv.FormatInt(ds.Result.Item.NumIid, 10),
		Title:      ds.Result.Item.Title,
		Price:      helper.PriceAsFloat(ds.Result.Item.PromotionPrice),
		ProductURL: ds.Result.Item.DetailURL,
		MainImgURL: ds.Result.Item.Images[0],
		Images:     ds.Result.Item.Images,
		DescImages: desc.GetImages(),
		Options:    sku.GetOptions(),
	}
}

func (c *DetailClient) getDesc(id string) model.Desc {
	ctx := context.TODO()
	s.Acquire(ctx, 1)
	defer s.Release(1)

	req := c.GetRequest(id, itemDesc)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var desc model.Desc

	json.Unmarshal(body, &desc)

	return desc
}

func (c *DetailClient) getSku(id string) model.Sku {
	ctx := context.TODO()
	s.Acquire(ctx, 1)
	defer s.Release(1)

	req := c.GetRequest(id, itemSku)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var sku model.Sku

	json.Unmarshal(body, &sku)

	return sku
}

func (c *DetailClient) getDetail(id string) model.DetailSimple {
	ctx := context.TODO()
	s.Acquire(ctx, 1)
	defer s.Release(1)

	req := c.GetRequest(id, detailSimple)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var detail model.DetailSimple

	json.Unmarshal(body, &detail)

	return detail
}
