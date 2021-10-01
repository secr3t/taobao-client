package client

import "fmt"

const defaultScheme = "https"

func GetUri(host, apiName, query string) string {
	return fmt.Sprintf("%s://%s/api?api=%s&%s", defaultScheme, host, apiName, query)
}
