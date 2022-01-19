package model

import (
	"encoding/json"
	"errors"
	"github.com/secr3t/taobao-client/helper"
	"github.com/secr3t/taobao-client/model"
	"github.com/tidwall/gjson"
	"math"
	"strconv"
	"strings"
)

type DetailBase struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *Data  `json:"data"`
}
type SkuProps struct {
	Name   string   `json:"prop_name"`
	Pid    string   `json:"pid"`
	Values []Values `json:"values"`
}
type Data struct {
	ItemID         int64      `json:"item_id"`
	Title          string     `json:"title"`
	CategoryID     string     `json:"category_id"`
	RootCategoryID string     `json:"root_category_id"`
	MainImgs       []string   `json:"main_imgs"`
	SkuProps       []SkuProps `json:"sku_props"`
	Skus           []BaseSku  `json:"skus"`
}
type BaseSku struct {
	SkuID       string  `json:"skuid"`
	SalePrice   *string `json:"sale_price"`
	OriginPrice *string `json:"origin_price"`
	Stock       string  `json:"stock"`
	PropPath    string  `json:"props_ids"`
}

func (d DetailBase) IsSuccess() bool {
	return d.Data != nil
}

func (d *Data) GetOptions(promotionRate float64) []model.Option {
	options := make([]model.Option, 0)
	priceMap := make(map[string]float64)
	optionMap := make(map[string]model.Option)

	for _, prop := range d.SkuProps {
		pid := prop.Pid
		name := prop.Name
		for _, value := range prop.Values {
			var img string
			if value.Image != nil {
				img = *value.Image
			}
			option := model.Option{
				Name:  name,
				Value: value.Name,
				Img:   img,
			}
			optionMap[pid+":"+value.Vid] = option
		}
	}

	for _, sku := range d.Skus {
		for _, propPath := range strings.Split(sku.PropPath, ";") {
			var price float64
			if sku.SalePrice == nil {
				if sku.OriginPrice == nil || promotionRate == 1 {
					continue
				}
				price = helper.PriceAsFloat(*sku.OriginPrice) * promotionRate
			} else {
				price = helper.PriceAsFloat(*sku.SalePrice)
			}
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

type Status struct {
	Msg           string  `json:"msg"`
	Code          int     `json:"code"`
	ExecutionTime float64 `json:"execution_time"`
}

func (s *Status) UnmarshalJSON(b []byte) error {
	r := gjson.ParseBytes(b)
	s.Msg = r.Get("msg").String()
	s.Code = int(r.Get("code").Int())
	s.ExecutionTime = r.Get("execution_time").Float()

	return nil
}

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
		if !strings.HasPrefix(img, "http") {
			img = "http:" + img
		}
		img = `<img src="` + img + `">`
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
	Image *string `json:"imageUrl"`
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
	Skus []Skus `json:"skus,omitempty"`
}

type SkuResult struct {
	Status  Status             `json:"status"`
	Item    *Item              `json:"item"`
	Prop    []Prop             `json:"prop,omitempty"`
	SkuMap  map[string]SkuInfo `json:"skus"`
	SkuBase SkuBase            `json:"sku_base"`
}

type SkuInfo struct {
	PromotionPrice float64 `json:"promotion_price"`
	Quantity       int     `json:"quantity"`
	Price          float64 `json:"price"`
}

func (s *SkuInfo) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var quantity string
	err = json.Unmarshal(*objMap["quantity"], &quantity)
	if err != nil {
		return err
	}

	s.Quantity, _ = strconv.Atoi(quantity)

	var price string
	if objMap["price"] != nil {
		err = json.Unmarshal(*objMap["price"], &price)
		if err != nil {
			var floatPrice float64
			if err = json.Unmarshal(*objMap["price"], &floatPrice); err != nil {
				return err
			} else {
				s.Price = floatPrice
			}
		} else {
			if strings.Contains(price, "-") {
				price = strings.Split(price, "-")[0]
			}

			s.Price, _ = strconv.ParseFloat(price, 64)
		}
	} else {
		return errors.New("price is nil")
	}

	if objMap["promotion_price"] == nil {
		s.PromotionPrice = s.Price
	} else {
		var promotionPrice string
		err = json.Unmarshal(*objMap["promotion_price"], &promotionPrice)
		if err != nil {
			return err
		}
		if strings.Contains(promotionPrice, "-") {
			promotionPrice = strings.Split(promotionPrice, "-")[0]
		}
		s.PromotionPrice, _ = strconv.ParseFloat(promotionPrice, 64)
	}

	return nil
}

func (s Sku) IsSuccess() bool {
	return s.Result.Status.Msg == "success" && s.Result.Item != nil
}

func (s Sku) GetOptions(promotionRate float64) []model.Option {
	if !s.IsSuccess() {
		return nil
	}
	options := make([]model.Option, 0)
	priceMap := make(map[string]float64)
	optionMap := make(map[string]model.Option)
	skuMap := make(map[string]string)

	for _, prop := range s.Result.Prop {
		pid := prop.Pid
		name := prop.Name
		for _, value := range prop.Values {
			var img string
			if value.Image != nil {
				img = "http:" + *value.Image
			}
			option := model.Option{
				Name:  name,
				Value: value.Name,
				Img:   img,
			}
			optionMap[pid+":"+value.Vid] = option
		}
	}

	for _, sku := range s.Result.SkuBase.Skus {
		skuMap[sku.SkuID] = sku.PropPath
	}

	for skuId, skuInfo := range s.Result.SkuMap {
		if skuId == "0" {
			continue
		}
		for _, propPath := range strings.Split(skuMap[skuId], ";") {
			price := skuInfo.PromotionPrice
			if price == skuInfo.Price && promotionRate != 1 {
				price = price * promotionRate
			}
			if price == 0 {
				price = skuInfo.Price * promotionRate
			}
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
