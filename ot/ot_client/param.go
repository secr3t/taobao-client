package client

import (
	"encoding/xml"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	cn = "zh-chs"
	en = "en"
)

var regex = regexp.MustCompile(`\p{Han}`)

type SearchParam struct {
	Page          int
	PageSize      int
	XmlParameters *SearchItemsParameters
}

type SearchItemsParameters struct {
	Provider           string    `xml:"Provider"`
	SearchMethod       string    `xml:"SearchMethod"`
	CurrencyCode       string    `xml:"CurrencyCode"`
	IsSellAllowed      bool      `xml:"IsSellAllowed"`
	UseOptimalFameSize bool      `xml:"UseOptimalFameSize"`
	Features           []Feature `xml:"Features"`
	LanguageOfQuery    *string   `xml:"LanguageOfQuery"`
	BrandId            *string   `xml:"BrandId,omitempty"` // ppath
	CategoryId         *string   `xml:"CategoryId,omitempty"`
	VendorName         *string   `xml:"VendorName,omitempty"`
	VendorId           *string   `xml:"VendorId,omitempty"`
	VendorAreaId       *string   `xml:"VendorAreaId,omitempty"`
	ItemTitle          *string   `xml:"ItemTitle,omitempty"` // q
	MinPrice           *float64  `xml:"MinPrice,omitempty"`  // start price
	MaxPrice           *float64  `xml:"MaxPrice,omitempty"`  // end price
}

type Feature struct {
	Name Name `xml:"Feature"`
}

type Name struct {
	Name   string `xml:"Name,attr"`
	Enable bool   `xml:",chardata"`
}

func NewParams() *SearchItemsParameters {
	return &SearchItemsParameters{
		Provider:           "Taobao",
		SearchMethod:       "Official",
		CurrencyCode:       "CNY",
		IsSellAllowed:      true,
		UseOptimalFameSize: true,
		Features: []Feature{
			{Name{"InStock", true}},
		},
	}
}

func (p *SearchItemsParameters) CatId(categoryId string) *SearchItemsParameters {
	if categoryId != "" {
		p.CategoryId = &categoryId
	}
	return p
}

func (p *SearchItemsParameters) VName(vendorName string) *SearchItemsParameters {
	p.VendorName = &vendorName
	return p
}

func (p *SearchItemsParameters) VId(vendorId string) *SearchItemsParameters {
	p.VendorId = &vendorId
	return p
}

func (p *SearchItemsParameters) VAreaId(vendorAreaId string) *SearchItemsParameters {
	p.VendorAreaId = &vendorAreaId
	return p
}

func (p *SearchItemsParameters) StartPrice(startPrice float64) *SearchItemsParameters {
	if startPrice != 0.0 {
		p.MinPrice = &startPrice
	}
	return p
}

func (p *SearchItemsParameters) EndPrice(endPrice float64) *SearchItemsParameters {
	if endPrice != 0.0 {
		p.MaxPrice = &endPrice
	}
	return p
}

func (p *SearchItemsParameters) Q(query string) *SearchItemsParameters {
	if hasChinese(query) {
		p.Cn()
	} else {
		p.En()
	}
	p.ItemTitle = &query
	return p
}

func (p *SearchItemsParameters) Cn() {
	loq := cn
	p.LanguageOfQuery = &loq
}

func (p *SearchItemsParameters) En() {
	loq := en
	p.LanguageOfQuery = &loq
}

func (p *SearchItemsParameters) Ppath(ppath string) *SearchItemsParameters {
	if ppath != "" {
		brandId := strings.Split(ppath, ":")[1]
		p.BrandId = &brandId
	}
	return p
}

func (p *SearchItemsParameters) ToXml() string {
	bytes, err := xml.Marshal(p)
	if err != nil {
		log.Print(err)
		return ""
	}

	return string(bytes)
}

func (p SearchParam) ToQuery(apiKey string) string {
	query := url.Values{}

	query.Add("instanceKey", apiKey)
	query.Add("frameSize", strconv.Itoa(p.PageSize))
	query.Add("framePosition", strconv.Itoa(p.Page))
	query.Add("xmlParameters", p.XmlParameters.ToXml())
	query.Add("blockList", "")

	return query.Encode()
}

func SearchParamFromUri(page int, uri string) SearchParam {
	parse, _ := url.Parse(uri)
	values := parse.Query()

	sp, ep := GetStartEndPrice(values.Get("filter"))

	p := NewParams().Q(values.Get("q")).
		//CatId(values.Get("cat")).
		Ppath(values.Get("ppath")).
		StartPrice(sp).
		EndPrice(ep)

	return SearchParam{
		Page:          page,
		PageSize:      40,
		XmlParameters: p,
	}
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

func hasChinese(before string) bool {
	return regex.MatchString(before)
}
