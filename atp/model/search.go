package model

import "strings"

const picUrlPrefix = "http"

type Item struct {
	Title          string `json:"title"`
	PicURL         string `json:"pic_url"`
	PromotionPrice string `json:"promotion_price"`
	Price          string `json:"price"`
	Sales          int    `json:"sales"`
	NumIid         string `json:"num_iid"`
	SellerNick     string `json:"seller_nick"`
	SellerID       int    `json:"seller_id"`
	DetailURL      string `json:"detail_url"`
}

func (i Item)GetPicURL() string {
	if strings.HasPrefix(i.PicURL, picUrlPrefix) {
		return i.PicURL
	}

	return picUrlPrefix + ":" + i.PicURL
}

type Items struct {
	Page             string `json:"page"`
	RealTotalResults int    `json:"real_total_results"`
	TotalResults     int    `json:"total_results"`
	PageSize         int    `json:"page_size"`
	PageCount        int    `json:"pagecount"`
	DataFrom         string `json:"data_from"`
	Item             []Item `json:"item"`
}

type SearchResult struct {
	Items     *Items  `json:"items"`
	ErrorCode *string `json:"error_code"`
	Reason    *string `json:"reason"`
	Error     string  `json:"error"`
}

func (r SearchResult) IsError() bool {
	return r.Error != ""
}
