package ruleengine

import (
	"github.com/curltech/go-colla-core/cache"
)

const (
	ExecuteType_GoEngine = "GoEngine"
	ExecuteType_GoRule   = "GoRule"
)

var MemCache = cache.NewMemCache("ruleengine", 1, 10)
