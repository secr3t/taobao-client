package model

import (
	"regexp"
	"strings"
)

const httpsPrefix = "https:"

var (
	imgRegex, _ = regexp.Compile(`//.*/imgextra/.*\.jpg`)
)

type Detail struct {
	Result struct {
		Item   DetailItem   `json:"item"`
		Status DetailStatus `json:"status"`
	} `json:"result"`
	RateLimit *RateLimit
}

type DetailItem struct {
	Title    string               `json:"title"`
	Images   []string             `json:"images"`
	DescImgs []string             `json:"desc_imgs"`
	NumIid   string               `json:"num_iid"`
	Skus     map[string]SkuDetail `json:"skus"`
	SkuBase  struct {
		Skus []struct {
			PropPath string `json:"propPath"`
			SkuID    string `json:"skuId"`
		} `json:"skus"`
		Prop []struct {
			Values []PropValue `json:"values"`
			Name   string      `json:"name"`
			Pid    string      `json:"pid"`
		} `json:"prop"`
	} `json:"sku_base"`
	DescUrl   string `json:"desc_url"`
	DetailUrl string `json:"detail_url"`
}

type PropValue struct {
	Vid   string `json:"vid"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type SkuDetail struct {
	PromotionPrice string `json:"promotion_price"`
	Quantity       string `json:"quantity"`
	Price          string `json:"price"`
}

type DetailStatus struct {
	Msg           string `json:"msg"`
	ExecutionTime string `json:"execution_time"`
	Code          int    `json:"code"`
}

func (d Detail) IsSuccess() bool {
	return d.Result.Status.Msg == "success"
}

func (i DetailItem) GetImages() []string {
	var images []string
	for _, imgUrl := range i.Images {
		if imgRegex.MatchString(imgUrl) {
			if strings.HasPrefix(imgUrl, httpsPrefix) {
				images = append(images, imgUrl)
			} else {
				images = append(images, httpsPrefix+imgUrl)
			}
		}
	}

	return images
}

func (i DetailItem) GetDetailUrl() string {
	return httpsPrefix + i.DetailUrl
}

func (i DetailItem) GetDescUrl() string {
	return httpsPrefix + i.DescUrl
}

func (i DetailItem) GetDescImgs() []string {
	var descImgs []string

	for _, descImg := range i.DescImgs {
		if strings.HasPrefix(descImg, "http") {
			descImgs = append(descImgs, descImg)
		} else {
			descImgs = append(descImgs, httpsPrefix+descImg)
		}
	}

	return descImgs
}

func (v PropValue) GetImage() string {
	if v.Image == "" {
		return ""
	}

	return httpsPrefix + v.Image
}
