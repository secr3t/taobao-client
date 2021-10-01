package model

import (
	model2 "github.com/secr3t/taobao-client/model"
	"regexp"
	"strings"
)

const httpsPrefix = "https:"

var (
	imgRegex, _ = regexp.Compile(`//.*/imgextra/.*\.jpg`)
)

type DetailItem struct {
	Title     string   `json:"title"`
	Images    []string `json:"images"`
	DescImgs  []string `json:"desc_imgs"`
	NumIid    string   `json:"num_iid"`
	DescUrl   string   `json:"desc_url"`
	DetailUrl string   `json:"detail_url"`
	Options   []model2.Option `json:"options"`
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
