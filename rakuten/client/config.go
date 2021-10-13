package client

import (
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
	keyRing *SafeRing
)

type ApiKey struct {
	Key    string
	Remain int
	Used   bool
}

func InitApiKeys(apiKeys ...string) {
	keyRing = NewSafeRing(apiKeys...)
}

func GetApiKey() string {
	return keyRing.Get()
}

type SafeRing struct {
	mutex *sync.Mutex
	strs  []string
	idx   int
	len   int
}

func NewSafeRing(strs ...string) *SafeRing {
	sr := &SafeRing{
		mutex: &sync.Mutex{},
		strs:  strs,
		idx:   0,
		len:   len(strs),
	}

	return sr
}

func (sr *SafeRing) Get() string {
	sr.mutex.Lock()
	defer func() {
		if sr.idx+1 == sr.len {
			sr.idx = 0
		} else {
			sr.idx += 1
		}
		sr.mutex.Unlock()
	}()

	return sr.strs[sr.idx]
}
