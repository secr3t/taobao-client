package client

import (
	"math"
	"sort"
	"sync"
)

const (
	taobaoApiHost = "taobao-api.p.rapidapi.com"
	searchItems   = "item_search"
	detailSimple  = "item_detail_simple"
	itemDesc      = "item_desc"
	itemSku       = "item_sku"
	httpsPrefix   = "https:"
)

var (
	keyMap = make(map[string]*ApiKey)
	keys   []*ApiKey
	m      = &sync.Mutex{}
)

type ApiKey struct {
	Key    string
	Remain int
	Used   bool
}

func InitApiKeys(apiKeys ...string) {
	for _, key := range apiKeys {
		apiKey := &ApiKey{Key: key, Remain: 0, Used: false}
		keyMap[key] = apiKey
		keys = append(keys, apiKey)
	}
}

func GetApiKey() string {
	m.Lock()
	defer m.Unlock()

	sort.Slice(keys, func(i, j int) bool {
		if !keys[i].Used && keys[j].Used {
			return true
		} else if keys[i].Used && !keys[j].Used {
			return false
		} else {
			return keys[i].Remain > keys[j].Remain
		}
	})

	return keys[0].Key
}

func ApiKeyUseEnd(apiKey string, remain int) {
	if !keyMap[apiKey].Used {
		keyMap[apiKey].Remain = remain
	} else {
		keyMap[apiKey].Remain = int(math.Min(float64(remain), float64(keyMap[apiKey].Remain)))
	}
	keyMap[apiKey].Used = true
}
