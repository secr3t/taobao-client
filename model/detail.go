package model

type DetailItem struct {
	Id             string
	Title          string
	Price          float64
	PromotionPrice float64
	ProductURL     string
	MainImgURL     string
	Images         []string
	DescImages     []string
	Options        []Option
}
