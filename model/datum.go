package model

type Datum struct {
	DescHtml     string   `json:"descHtml"`
	DescUrl      string   `json:"descUrl"`
	ItemId       string   `json:"itemId"`
	MainImg      string   `json:"mainImg"`
	Imgs         []string `json:"imgs"`
	MainImgUrl   string   `json:"mainImgUrl"`
	Name         string   `json:"name"`
	Options      []Option `json:"options"`
	PriceInChina float64  `json:"priceInChina"`
	ProductUrl   string   `json:"productUrl"`
}

type Option struct {
	Img   string  `json:"img"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Value string  `json:"value"`
}

type DatumDto struct {
	Data []Datum `json:"data"`
	Path string  `json:"path"`
}
