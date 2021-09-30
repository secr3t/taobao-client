package atp_client

import (
	"fmt"
	"sync/atomic"
)

const (
	host1  = "https://asia.atphosting24.com/taobao/index.php"
	host2  = "https://laonet.online/index.php"
	route = "api_tester/call"
	lang  = "zh-CN"
	sort  = "sale"
)

var (
	vv atomic.Value
	hosts = []string{host1, host2}
)

func init() {
	vv.Store(0)
}

func GetUri(query string) string {
	return fmt.Sprintf("%s?%s", getHost(), query)
}

func getHost() string {
	idx := vv.Load().(int)
	if idx == 0 {
		vv.Swap(1)
	} else {
		vv.Swap(0)
	}
	return hosts[idx]
}
