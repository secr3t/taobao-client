package model

type Search struct {
	Result struct {
		Item         []SearchItem `json:"item"`
		TotalResults int          `json:"total_results"`
		Status       SearchStatus `json:"status"`
		PageSize     string       `json:"page_size"`
	} `json:"result"`
	RateLimit *RateLimit
}

type SearchItem struct {
	PromotionPrice string `json:"promotion_price"`
	Loc            string `json:"loc"`
	NumIid         int64  `json:"num_iid"`
	Usertype       int    `json:"usertype"`
	Pic            string `json:"pic"`
	Title          string `json:"title"`
	Sales          int    `json:"sales"`
	SellerId       int64  `json:"seller_id"`
	SellerNick     string `json:"seller_nick"`
	DeliveryFee    string `json:"delivery_fee"`
	Price          string `json:"price"`
	DetailUrl      string `json:"detail_url"`
	ShopTitle      string `json:"shop_title"`
}

type SearchStatus struct {
	Msg           string `json:"msg"`
	ExecutionTime string `json:"execution_time"`
	Code          int    `json:"code"`
	Action        int    `json:"action"`
}

func (d Search) IsSuccess() bool {
	return d.Result.Status.Msg == "success"
}
