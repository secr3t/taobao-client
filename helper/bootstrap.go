package helper

import (
	"net/url"
	"regexp"
	"strconv"
)

func ParseUri(uri string) url.Values {
	parse, _ := url.Parse(uri)
	return parse.Query()
}

func GetStartEndPrice(filter string) (startPrice int, endPrice int) {
	reg, _ := regexp.Compile(`reserve_price\[(\d+)?,(\d+)?\]`)

	matched := reg.FindStringSubmatch(filter)

	if len(matched) > 2 {
		startPrice, _ = strconv.Atoi(matched[1])
		endPrice, _ = strconv.Atoi(matched[2])
	}

	if len(matched) == 2 {
		startPrice, _ = strconv.Atoi(matched[1])
	}

	return
}

func PriceAsFloat(price string) float64 {
	p, _ := strconv.ParseFloat(price, 64)
	return p
}
