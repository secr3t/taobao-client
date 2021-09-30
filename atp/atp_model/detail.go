package atp_model

import (
	"github.com/secr3t/taobao-client/model"
	"math"
	"strconv"
	"strings"
)

type DetailResult struct {
	DetailItem *DetailItem `json:"item"`
	Error      string      `json:"error"`
}

func (dr DetailResult) IsError() bool {
	return dr.Error != ""
}

type ItemImgs struct {
	Url string `json:"url"`
}

type Sku struct {
	Price          string `json:"price"`
	TotalPrice     int    `json:"total_price"`
	OrginalPrice   string `json:"orginal_price"`
	Properties     string `json:"properties"`
	PropertiesName string `json:"properties_name"`
	Quantity       string `json:"quantity"`
	SkuID          string `json:"sku_id"`
}

type Skus struct {
	Sku []Sku `json:"sku"`
}

type DetailItem struct {
	NumIid    string            `json:"num_iid"`
	Title     string            `json:"title"`
	Price     string            `json:"price"`
	DetailURL string            `json:"detail_url"`
	PicURL    string            `json:"pic_url"`
	Desc      string            `json:"desc"`
	ItemImgs  []ItemImgs        `json:"item_imgs"`
	Skus      Skus              `json:"skus"`
	PropsList map[string]string `json:"props_list"`
	PropsImg  map[string]string `json:"props_img"`
	DescImg   []string          `json:"desc_img"`
	Options   []model.Option
}

func (di *DetailItem) SetOptions() {
	options := make([]model.Option, 0)

	priceMap := make(map[string]float64)

	for _, sku := range di.Skus.Sku {
		for _, propPath := range strings.Split(sku.Properties, ";") {
			price, _ := strconv.ParseFloat(sku.Price, 64)
			if val, ok := priceMap[propPath]; ok {
				priceMap[propPath] = math.Min(val, price)
			} else {
				priceMap[propPath] = price
			}
		}
	}

	for propPath, value := range di.PropsList {
		splited := strings.Split(value, ":")
		option := model.Option{
			Name:     splited[0],
			Value:    splited[1],
			Img:      di.GetPropImg(propPath),
			Price:    priceMap[propPath],
		}
		options = append(options, option)
	}
	di.Options = options
}

func (di *DetailItem) GetMainImg() string {
	return ValidImg(di.PicURL)
}

func (di *DetailItem) GetItemImgs() []string {
	imgs := make([]string, 0)
	for _, itemImgs := range di.ItemImgs {
		img := itemImgs.Url
		imgs = append(imgs, ValidImg(img))
	}

	return imgs
}

func (di *DetailItem) GetPropImg(propPath string) string {
	img := di.PropsImg[propPath]

	return ValidImg(img)
}

func ValidImg(img string) string {
	if img == "" {
		return ""
	}
	if !strings.HasPrefix(img, "http") {
		img = "http:" + img
	}
	return img
}
