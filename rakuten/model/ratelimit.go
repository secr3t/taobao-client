package model

import (
	"net/http"
	"strconv"
)

const (
	limitHeader = "x-ratelimit-optimal-limit"
	remainHeader = "x-ratelimit-optimal-remaining"
	resetHeader = "x-ratelimit-optimal-reset"
)

type RateLimit struct {
	Limit    int
	Remain   int
	ResetSec int
}

func FromHeader(header http.Header) *RateLimit {
	limit, _ := strconv.Atoi(header.Get(limitHeader))
	remain, _ := strconv.Atoi(header.Get(remainHeader))
	resetSec, _ := strconv.Atoi(header.Get(resetHeader))
	return &RateLimit{
		limit,
		remain,
		resetSec,
	}
}
