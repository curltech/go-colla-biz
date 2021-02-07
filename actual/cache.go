package actual

import (
	"fmt"
	"github.com/curltech/go-colla-core/logger"
)

func getKey(schemaName string, id uint64) string {
	return fmt.Sprintf("Role:%v:%v", schemaName, id)
}

func getCacheRole(schemaName string, id uint64) *Role {
	key := getKey(schemaName, id)
	v, ok := MemCache.Get(key)
	if ok {
		role := v.(*Role)

		return role
	}

	return nil
}

func setCacheRole(role *Role) {
	key := getKey(role.SchemaName, role.Id)
	_, ok := MemCache.Get(key)
	if ok {
		logger.Sugar.Errorf("MemCacheExist:%v", role.Id)
	}
	MemCache.SetDefault(key, role)
}

func removeCacheRole(schemaName string, id uint64) {
	key := getKey(schemaName, id)
	MemCache.Delete(key)
}
