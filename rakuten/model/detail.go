package model

import (
	"github.com/secr3t/taobao-client/model"
	"math"
	"strconv"
	"strings"
)

// DetailSimple struct start
type DetailSimple struct {
	Result Result `json:"result"`
}
type Status struct {
	Msg           string  `json:"msg"`
	Code          int     `json:"code"`
	ExecutionTime float64 `json:"execution_time"`
}
type DetailItem struct {
	NumIid         int64    `json:"num_iid"`
	Title          string   `json:"title"`
	TotalSales     int      `json:"total_sales"`
	DetailURL      string   `json:"detail_url"`
	Images         []string `json:"images"`
	PromotionPrice string   `json:"promotion_price"`
	Price          string   `json:"price"`
}
type Result struct {
	Status Status      `json:"status"`
	Item   *DetailItem `json:"item"`
}

func (d DetailSimple) IsSuccess() bool {
	return d.Result.Status.Msg == "success"
}
// DetailSimple struct end

// Desc struct start
type Desc struct {
	Result DescResult `json:"result"`
}

type DescResult struct {
	Status Status   `json:"status"`
	Item   []string `json:"item"`
}

func (d Desc) IsSuccess() bool {
	return d.Result.Status.Msg == "success" && d.Result.Item != nil
}

func (d Desc) GetImages() []string {
	var imgs []string
	for _, img := range d.Result.Item {
		if strings.HasPrefix(img, "http") {
			img = "http://" + img
		}
		imgs = append(imgs, img)
	}
	return imgs
}
// Desc struct end

// Sku struct start
type Sku struct {
	Result SkuResult `json:"result"`
}
type Item struct {
	Pic            string `json:"pic"`
	Price          string `json:"price"`
	PromotionPrice string `json:"promotion_price"`
	Quantity       string `json:"quantity"`
}
type Values struct {
	Vid   string  `json:"vid"`
	Name  string  `json:"name"`
	Image *string `json:"image"`
}
type Prop struct {
	Pid    string   `json:"pid"`
	Name   string   `json:"name"`
	Values []Values `json:"values"`
}
type Skus struct {
	SkuID    string `json:"skuId"`
	PropPath string `json:"propPath"`
}
type SkuBase struct {
	Skus []Skus `json:"skus"`
}
type SkuResult struct {
	Status  Status             `json:"status"`
	Item    *Item              `json:"item"`
	Prop    []Prop             `json:"prop"`
	SkuMap  map[string]SkuInfo `json:"skus"`
	SkuBase *SkuBase           `json:"sku_base"`
}

type SkuInfo struct {
	PromotionPrice string `json:"promotion_price"`
	Quantity       string `json:"quantity"`
	Price          string `json:"price"`
}

func (d Sku) IsSuccess() bool {
	return d.Result.Status.Msg == "success" && d.Result.Item != nil
}

func (d Sku) GetOptions() []model.Option {
	if !d.IsSuccess() {
		return nil
	}
	options := make([]model.Option, 0)
	priceMap := make(map[string]float64)
	optionMap := make(map[string]model.Option)
	skuMap := make(map[string]string)

	for _, prop := range d.Result.Prop {
		pid := prop.Pid
		name := prop.Name
		for _, value := range prop.Values {
			var img string
			if value.Image != nil {
				img = "http:" + *value.Image
			}
			option := model.Option{
				Name: name,
				Value: value.Name,
				Img: img,
			}
			optionMap[pid + ":" + value.Vid] = option
		}
	}

	for _, sku := range d.Result.SkuBase.Skus {
		skuMap[sku.SkuID] = sku.PropPath
	}

	for skuId, skuInfo := range d.Result.SkuMap {
		if skuId == "0" {
			continue
		}
		for _, propPath := range strings.Split(skuMap[skuId], ";") {
			price, _ := strconv.ParseFloat(skuInfo.PromotionPrice, 64)
			if val, ok := priceMap[propPath]; ok {
				priceMap[propPath] = math.Min(val, price)
			} else {
				priceMap[propPath] = price
			}
		}
	}

	for propPath, option := range optionMap {
		option.Price = priceMap[propPath]
		options = append(options, option)
	}

	return options
}
// Sku struct end
