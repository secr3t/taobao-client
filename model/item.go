package model

type Item struct {
	Id         string
	Title      string
	CategoryId string
	ProductUrl string
	MainImgUrl string
	Price      float64
	Imgs       []string
}

type SearchResult struct {
	TotalCount int
	FrameSize  int
	Items      []Item
}
