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
	"sync"
)

type SearchClient struct {
	ApiKey string
}

func NewSearchClient(apiKey string) *SearchClient {
	return &SearchClient{
		ApiKey: apiKey,
	}
}

func (c *SearchClient) GetItems(uri string) []model2.Item {

	_, e := url.ParseRequestURI(uri)

	if e != nil {
		return []model2.Item{}
	}

	rakutenParam := ParamsFromUri(uri)
	return c.SearchTilLimit(&rakutenParam, 200)
}

func (c *SearchClient) SearchTilLimit(param *SearchParam, limit int) []model2.Item {
	result := c.SearchItems(*param)

	if limit > result.Result.TotalResults {
		limit = result.Result.TotalResults
	}

	itemsChain := make(chan model.SearchItem, limit)

	limit -= param.PageSize

	for _, item := range result.Result.Item {
		itemsChain <- item
	}

	wg := sync.WaitGroup{}

	for ; limit > 0; limit -= param.PageSize {
		param.Page += 1
		wg.Add(1)
		go func() {
			result = c.SearchItems(*param)
			for _, item := range result.Result.Item {
				itemsChain <- item
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(itemsChain)
	}()

	items := make([]model2.Item, 0)
	for rakutenSearchItem := range itemsChain {
		price, _ := strconv.ParseFloat(rakutenSearchItem.Price, 64)
		item := model2.Item{
			Id:         strconv.FormatInt(rakutenSearchItem.NumIid, 10),
			Title:      rakutenSearchItem.Title,
			ProductUrl: rakutenSearchItem.DetailUrl,
			MainImgUrl: rakutenSearchItem.Pic,
			Price:      price,
		}
		items = append(items, item)
	}

	return items
}

func (c *SearchClient) SearchItems(param SearchParam) model.Search {
	query := param.ToQuery()

	uri := GetUri(taobaoApiHost, searchItems, query)

	req, _ := http.NewRequest("GET", uri, nil)

	req.Header.Add("x-rapidapi-key", c.ApiKey)
	req.Header.Add("x-rapidapi-host", taobaoApiHost)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var search model.Search

	json.Unmarshal(body, &search)

	rateLimit := model.FromHeader(res.Header)
	search.RateLimit = rateLimit

	return search
}

func ParamsFromUri(uri string) SearchParam {
	values := helper.ParseUri(uri)

	catId, _ := strconv.Atoi(values.Get("cat"))
	startPrice, endPrice := helper.GetStartEndPrice(values.Get("filter"))

	return NewSearchParam(
		values.Get("q"),
		values.Get("sort"),
		1,
		100,
		startPrice,
		endPrice,
		catId,
	)
}
