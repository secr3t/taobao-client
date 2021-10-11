package model

import (
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
// Sku struct end
