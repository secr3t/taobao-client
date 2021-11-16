package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/secr3t/taobao-client/helper"
	model2 "github.com/secr3t/taobao-client/model"
	otClient "github.com/secr3t/taobao-client/ot/client"
	"github.com/secr3t/taobao-client/rakuten/model"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type DetailClient struct {
	OtClient *otClient.DetailClient
	s        *semaphore.Weighted
	keyRing  *SafeRing
	hook func(key string)
}

func NewDetailClient(weight int64, keys []string, hook func(key string)) *DetailClient {
	return &DetailClient{
		s: semaphore.NewWeighted(weight),
		keyRing: NewSafeRing(keys...),
		hook: hook,
	}
}

func (c *DetailClient) AddOtClient(otApiKey string) *DetailClient {
	c.OtClient = otClient.NewDetailClient(otApiKey)
	return c
}

func (c *DetailClient) GetRequest(id, api string) (*http.Request, string) {
	queryV := url.Values{}
	queryV.Add("num_iid", id)
	uri := GetUri(taobaoApiHost, api, queryV.Encode())

	req, _ := http.NewRequest("GET", uri, nil)

	key := c.keyRing.Get()
	req.Header.Add("x-rapidapi-key", key)
	req.Header.Add("x-rapidapi-host", taobaoApiHost)

	if c.hook != nil {
		go c.hook(key)
	}

	return req, key
}

func (c *DetailClient) GetDetail(id string) (*model2.DetailItem, error) {
	var detailItem *model2.DetailItem
	ds := c.getDetail(id)

	if !ds.IsSuccess() {
		if c.OtClient == nil {
			return nil, errors.New("detail : rakuten fail, ot empty " + id)
		}
		if detailItem = c.OtClient.GetDetailBase(id); detailItem == nil {
			return nil, errors.New("detail : rakuten fail, ot fail " + id)
		}
	}

	desc := c.getDesc(id)

	if !desc.IsSuccess() {
		return nil, errors.New("desc : rakuten fail, " + id)
	}

	sku := c.getSku(id)

	if !sku.IsSuccess() {
		return nil, errors.New("sku : rakuten fail, " + id)
	}

	if detailItem != nil {
		options := sku.GetOptions()
		var price float64
		for _, option := range options {
			if price == 0 {
				price = option.Price
			} else if option.Price != 0{
				price = math.Min(price, option.Price)
			}
		}
		detailItem.Options = options
		detailItem.Price = price
		detailItem.DescImages = desc.GetImages()

		return detailItem, nil
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
	}, nil
}

func (c *DetailClient) getDesc(id string) model.Desc {
	ctx := context.TODO()
	c.s.Acquire(ctx, 1)
	defer c.s.Release(1)

	req, _ := c.GetRequest(id, itemDesc)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return model.Desc{}
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var desc model.Desc

	json.Unmarshal(body, &desc)

	return desc
}

func (c *DetailClient) getSku(id string) model.Sku {
	ctx := context.TODO()
	c.s.Acquire(ctx, 1)
	defer c.s.Release(1)

	req, _ := c.GetRequest(id, itemSku)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var sku model.Sku

	err := json.Unmarshal(body, &sku)

	r := gjson.ParseBytes(body)

	if err != nil && !r.Get("result.skus").IsArray() {
		log.Println(err, id, res.Header)
	}

	return sku
}

func (c *DetailClient) getDetail(id string) model.DetailSimple {
	ctx := context.TODO()
	c.s.Acquire(ctx, 1)
	defer c.s.Release(1)

	req, _ := c.GetRequest(id, detailSimple)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var detail model.DetailSimple

	json.Unmarshal(body, &detail)

	return detail
}
