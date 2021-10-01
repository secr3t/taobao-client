package client

import "fmt"

const defaultScheme = "http"

func GetUri(host, apiName, query string) string {
	return fmt.Sprintf("%s://%s/%s?%s", defaultScheme, host, apiName, query)
}
