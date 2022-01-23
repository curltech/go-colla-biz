package ruleengine

import (
	"github.com/curltech/go-colla-core/cache"
)

const (
	ExecuteType_GoEngine = "GoEngine"
	ExecuteType_GoRule   = "GoRule"
	ExecuteType_TenGo    = "TenGo"
)

var MemCache = cache.NewMemCache("ruleengine", 60, 10)
