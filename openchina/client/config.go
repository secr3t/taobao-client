package client

import (
	"fmt"
)

const (
	host  = "https://api.openchinaapi.com/v1/taobao/products/"
)

func GetUri(itemId string) string {
	return fmt.Sprintf("%s%s?cache=no", host, itemId)
}

