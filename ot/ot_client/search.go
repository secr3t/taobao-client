package client

import (
	"errors"
	model "github.com/secr3t/taobao-client/ot/ot_model"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const (
	host        = "otapi.net/service-json"
	searchItems = "BatchSearchItemsFrame"

	IdPath           = "Id"
	TitlePath        = "OriginalTitle"
	CategoryIdPath   = "ExternalCategoryId"
	ProductUrlPath   = "ExternalItemUrl"
	MainImgUrlPath   = "MainPictureUrl"
	PricePath        = "Price.OriginalPrice"
	ImgsPath         = "Pictures.#.Url"
	totalCountPath   = "Result.Items.Items.TotalCount"
	currentFrameSize = "Result.Items.CurrentFrameSize"
	ErrorPartialUrl  = "img.alicdn.com/imgextra///img.alicdn.com/imgextra"
	FixPartialUrl    = "img.alicdn.com/imgextra"
)

var (
	SearchFail = errors.New("search failed")
)

type SearchClient struct {
	ApiKey string
}

func NewSearchClient(apiKey string) *SearchClient {
	return &SearchClient{
		ApiKey: apiKey,
	}
}

func (c *SearchClient) SearchItems(param SearchParam) (model.SearchResult, error) {
	query := param.ToQuery(c.ApiKey)

	uri := GetUri(host, searchItems, query)

	req, _ := http.NewRequest("GET", uri, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return searchResultToItems(body)
}

func (c *SearchClient) SearchTilLimit(param *SearchParam, limit int)[]model.Item  {
	result, err := c.SearchItems(*param)

	if err != nil {
		return nil
	}

	if limit > result.TotalCount {
		limit = result.TotalCount
	}

	itemsChan := make(chan model.Item, limit)

	frameSize := result.FrameSize
	limit -= frameSize

	for _, item := range result.Items {
		itemsChan <- item
	}

	wg := sync.WaitGroup{}

	for ;limit > 0; limit -= frameSize{
		param.Page += len(result.Items)
		wg.Add(1)
		go func() {
			result, err = c.SearchItems(*param)
			for _,item := range result.Items {
				itemsChan <- item
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

func searchResultToItems(json []byte) (model.SearchResult, error) {
	jsonStr := strings.Replace(string(json), ErrorPartialUrl, FixPartialUrl, -1)
	r := gjson.Parse(jsonStr)

	if r.Get("ErrorCode").String() != "Ok" {
		return model.SearchResult{}, SearchFail
	}

	var items = make([]model.Item, 0)

	r.Get("Result.Items.Items.Content").ForEach(func(key, value gjson.Result) bool {
		items = append(items, model.Item{
			Id:         value.Get(IdPath).String(),
			Title:      value.Get(TitlePath).String(),
			CategoryId: value.Get(CategoryIdPath).String(),
			ProductUrl: value.Get(ProductUrlPath).String(),
			MainImgUrl: value.Get(MainImgUrlPath).String(),
			Price:      value.Get(PricePath).Float(),
			Imgs:       ConvertImgUrls(value.Get(ImgsPath).Array()),
		})
		return true
	})

	var result model.SearchResult
	result.Items = items
	result.FrameSize = int(r.Get(currentFrameSize).Int())
	result.TotalCount = int(r.Get(totalCountPath).Int())

	return result, nil
}

func ConvertImgUrls(imgUrls []gjson.Result) []string {
	var imgs = make([]string, 0)

	for _, url := range imgUrls {
		imgs = append(imgs, url.String())
	}

	return imgs
}
