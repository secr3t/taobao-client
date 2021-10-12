package client

import (
	"encoding/json"
	"github.com/secr3t/taobao-client/helper"
	model2 "github.com/secr3t/taobao-client/model"
	"github.com/secr3t/taobao-client/rakuten/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
)

var idx int64

type DetailClient struct {
	ApiKeys []string
}

func NewDetailClient(apiKeys ...string) *DetailClient {
	return &DetailClient{
		ApiKeys: apiKeys,
	}
}

func (c *DetailClient) GetRequest(id, api string) *http.Request {
	queryV := url.Values{}
	queryV.Add("num_iid", id)
	uri := GetUri(taobaoApiHost, api, queryV.Encode())

	req, _ := http.NewRequest("GET", uri, nil)

	req.Header.Add("x-rapidapi-key", c.ApiKeys[idx])
	req.Header.Add("x-rapidapi-host", taobaoApiHost)

	if idx < int64(len(c.ApiKeys) - 1) {
		atomic.AddInt64(&idx, 1)
	} else {
		atomic.StoreInt64(&idx, 0)
	}

	return req
}

func (c *DetailClient) GetDetail(id string) *model2.DetailItem {
	ds := c.getDetail(id)

	if !ds.IsSuccess() {
		return nil
	}

	desc := c.getDesc(id)

	if !desc.IsSuccess() {
		return nil
	}

	sku := c.getSku(id)

	if !sku.IsSuccess() {
		return nil
	}

	return &model2.DetailItem{
		Id: strconv.FormatInt(ds.Result.Item.NumIid, 10),
		Title: ds.Result.Item.Title,
		Price: helper.PriceAsFloat(ds.Result.Item.PromotionPrice),
		ProductURL: ds.Result.Item.DetailURL,
		MainImgURL: ds.Result.Item.Images[0],
		Images: ds.Result.Item.Images,
		DescImages: desc.GetImages(),
		Options: sku.GetOptions(),
	}
}

func (c *DetailClient) getDesc(id string) model.Desc {
	req := c.GetRequest(id, itemDesc)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var desc model.Desc

	json.Unmarshal(body, &desc)

	return desc
}

func (c *DetailClient) getSku(id string) model.Sku {
	req := c.GetRequest(id, itemSku)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var sku model.Sku

	json.Unmarshal(body, &sku)

	return sku
}

func (c *DetailClient) getDetail(id string) model.DetailSimple {
	req := c.GetRequest(id, detailSimple)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var detail model.DetailSimple

	json.Unmarshal(body, &detail)

	return detail
}
