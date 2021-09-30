package atp_client

import (
	"encoding/json"
	"fmt"
	model "github.com/secr3t/taobao-client/atp/atp_model"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
)

const (
	searchApiName = "item_search"
)

type SearchClient struct {
	apiKey string
}

func NewSearchClient(apiKey string) *SearchClient {
	return &SearchClient{
		apiKey: apiKey,
	}
}

func (c *SearchClient) searchItems(param SearchParam) (model.SearchResult, error) {
	query := param.ToQueryParam()

	reqUri := GetUri(query)

	res, _ := http.Get(reqUri)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var result model.SearchResult

	err := json.Unmarshal(body, &result)

	if err != nil {
		return model.SearchResult{}, err
	}

	return result, nil
}

func (c *SearchClient) SearchItems(uri string) (model.SearchResult, error) {
	param := SearchParamFromUri(uri, c.apiKey)
	return c.searchItems(*param)
}

func (c *SearchClient) SearchTilLimit(uri string, limit int) []model.Item {
	param := SearchParamFromUri(uri, c.apiKey)
	result, err := c.searchItems(*param)

	if err != nil {
		return nil
	}

	if limit > result.Items.RealTotalResults {
		limit = result.Items.RealTotalResults
	}

	itemsChan := make(chan model.Item, limit)
	length := len(result.Items.Item)
	limit -= length

	for _, item := range result.Items.Item {
		itemsChan <- item
	}

	wg := sync.WaitGroup{}

	for ; limit > 0; limit -= length {
		param.page++
		searchParam := *param
		wg.Add(1)
		go func() {
			result, err = c.searchItems(searchParam)
			if result.Items != nil {
				for _, item := range result.Items.Item {
					itemsChan <- item
				}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(itemsChan)
	}()

	items := make([]model.Item, 0)
	for item := range itemsChan {
		items = append(items, item)
	}

	return items
}

type SearchParam struct {
	q          string
	startPrice float64
	endPrice   float64
	page       int
	pageSize   int
	ppath      string
	lang       string
	sort       string
	cat        string
	key        string
}

func SearchParamFromUri(uri, apiKey string) *SearchParam {
	parse, _ := url.Parse(uri)
	values := parse.Query()

	sp, ep := GetStartEndPrice(values.Get("filter"))

	return &SearchParam{
		q:          values.Get("q"),
		startPrice: sp,
		endPrice:   ep,
		page:       1,
		pageSize:   40,
		ppath:      values.Get("ppath"),
		lang:       lang,
		sort:       sort,
		cat:        values.Get("cat"),
		key:        apiKey,
	}
}

func (sp SearchParam) ToQueryParam() string {
	queryParams := url.Values{}

	queryParams.Add("api_name", searchApiName)
	queryParams.Add("route", route)
	queryParams.Add("q", sp.q)
	queryParams.Add("start_price", fmt.Sprint(sp.startPrice))
	queryParams.Add("end_price", fmt.Sprint(sp.endPrice))
	queryParams.Add("page_size", fmt.Sprint(sp.pageSize))
	queryParams.Add("page", fmt.Sprint(sp.page))
	queryParams.Add("ppath", sp.ppath)
	queryParams.Add("lang", sp.lang)
	queryParams.Add("sort", sp.sort)
	queryParams.Add("cat", sp.cat)
	queryParams.Add("key", sp.key)

	return queryParams.Encode()
}

func GetStartEndPrice(filter string) (startPrice float64, endPrice float64) {
	reg, _ := regexp.Compile(`reserve_price\[(\d*\.?\d*)?,(\d*\.?\d*)?\]`)

	matched := reg.FindStringSubmatch(filter)

	if len(matched) > 2 {
		startPrice, _ = strconv.ParseFloat(matched[1], 64)
		endPrice, _ = strconv.ParseFloat(matched[2], 64)
	}

	if len(matched) == 2 {
		startPrice, _ = strconv.ParseFloat(matched[1], 64)
	}

	return
}
