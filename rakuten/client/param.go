package client

import (
	"net/url"
	"strconv"
)

type SearchParam struct {
	Q          string
	Page       int
	StartPrice int
	EndPrice   int
	Sort       string
	PageSize   int
	CatId      int
}

func NewSearchParam(q, sort string, page, pageSize, startPrice, endPrice, catId int) SearchParam {
	return SearchParam{
		Q:          q,
		Page:       page,
		StartPrice: startPrice,
		EndPrice:   endPrice,
		Sort:       sort,
		PageSize:   pageSize,
		CatId:      catId,
	}
}

func (p SearchParam) ToQuery() string {
	query := url.Values{}

	query.Add("q", p.Q)

	if p.PageSize == 0 {
		p.PageSize = 40
	}

	if p.PageSize > 100 {
		p.PageSize = 100
	}

	query.Add("page_size", strconv.Itoa(p.PageSize))

	if p.Page != 0 {
		query.Add("page", strconv.Itoa(p.Page))
	}

	if p.StartPrice != 0 {
		query.Add("start_price", strconv.Itoa(p.StartPrice))
	}

	if p.EndPrice != 0 {
		query.Add("end_price", strconv.Itoa(p.EndPrice))
	}

	if p.CatId != 0 {
		query.Add("cat", strconv.Itoa(p.CatId))
	}

	if p.Sort != "" {
		query.Add("sort", p.Sort)
	}

	return query.Encode()
}
