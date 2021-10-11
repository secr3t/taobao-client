package client

import (
	"encoding/json"
	"github.com/secr3t/taobao-client/rakuten/model"
	"io/ioutil"
	"net/http"
	"net/url"
)

type DetailClient struct {
	ApiKey string
}

func NewDetailClient(apiKey string) *DetailClient {
	return &DetailClient{
		ApiKey: apiKey,
	}
}

func (c *DetailClient) GetRequest(id, api string) *http.Request {
	queryV := url.Values{}
	queryV.Add("num_iid", id)
	uri := GetUri(taobaoApiHost, api, queryV.Encode())

	req, _ := http.NewRequest("GET", uri, nil)

	req.Header.Add("x-rapidapi-key", c.ApiKey)
	req.Header.Add("x-rapidapi-host", taobaoApiHost)

	return req
}

func (c *DetailClient) GetDesc(id string) model.Desc {
	req := c.GetRequest(id, itemDesc)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var desc model.Desc

	json.Unmarshal(body, &desc)

	return desc
}

func (c *DetailClient) GetSku(id string) model.Sku {
	req := c.GetRequest(id, itemSku)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var sku model.Sku

	json.Unmarshal(body, &sku)

	return sku
}

func (c *DetailClient) GetDetail(id string) model.DetailSimple {
	req := c.GetRequest(id, detailSimple)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var detail model.DetailSimple

	json.Unmarshal(body, &detail)

	return detail
}
