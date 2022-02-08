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
	"strings"
)

type DetailClient struct {
	OtClient      *otClient.DetailClient
	s             *semaphore.Weighted
	keyRing       *SafeRing
	hook          func(key, name string)
	detailBaseKey string
}

func NewDetailClient(weight int64, keys []string, detailBaseKey string, hook func(key, name string)) *DetailClient {
	return &DetailClient{
		s:             semaphore.NewWeighted(weight),
		keyRing:       NewSafeRing(keys...),
		hook:          hook,
		detailBaseKey: detailBaseKey,
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
		go c.hook(key, "taobao-api")
	}

	return req, key
}

func (c *DetailClient) GetDetail(arg string) (*model2.DetailItem, error) {
	var promotionRate = 1.0
	var promotionPrice = 0.0
	args := strings.Split(arg, "-")
	id := args[0]

	if len(args) > 1 {
		promotionPrice = helper.PriceAsFloat(args[1])
	}

	base := c.getDetail(id)

	if !base.IsSuccess() {
		return nil, errors.New("detail : ttpd failed (no promo) " + id)
	}

	desc := c.GetDesc(id)

	if !desc.IsSuccess() {
		return nil, errors.New("desc : taobao-api failed, " + id)
	}

	options := base.Data.GetOptions(promotionRate)
	if promotionPrice == 0.0 {
		for _, option := range options {
			if promotionPrice == 0 {
				promotionPrice = option.Price
			} else if option.Price != 0 {
				promotionPrice = math.Min(promotionPrice, option.Price)
			}
		}
	}

	return &model2.DetailItem{
		Id:         strconv.FormatInt(base.Data.ItemID, 10),
		Title:      base.Data.Title,
		Price:      promotionPrice,
		MainImgURL: base.Data.MainImgs[0],
		Images:     base.Data.MainImgs,
		DescImages: desc.GetImages(),
		Options:    options,
	}, nil
}

func (c *DetailClient) GetDesc(id string) model.Desc {
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

	if err != nil && r.Get("result.skus").IsArray() {
		log.Println(err, id, res.Header)
	}

	return sku
}

func (c *DetailClient) getDetail(id string) model.DetailBase {
	ctx := context.TODO()
	c.s.Acquire(ctx, 1)
	defer c.s.Release(1)

	if c.hook != nil {
		go c.hook(c.detailBaseKey, "ttpd")
	}

	req := c.getDetailBaseRequest(id)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var detail model.DetailBase

	json.Unmarshal(body, &detail)

	return detail
}

func (c *DetailClient) getDetailBaseRequest(id string) *http.Request {
	ttpdUrl := "https://taobao-tmall-product-data.p.rapidapi.com/api/sc/taobao/item_detail?item_id=" + id

	req, _ := http.NewRequest("GET", ttpdUrl, nil)

	req.Header.Add("x-rapidapi-key", c.detailBaseKey)
	req.Header.Add("x-rapidapi-host", "taobao-tmall-product-data.p.rapidapi.com")

	return req
}
