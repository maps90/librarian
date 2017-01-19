package cache

import (
	"time"
	"fmt"
)

type CStorage interface {
	Set(key string, data string, expr time.Duration) error
	Get(key string) string
}

const (
	REDIS = "redis"
	MEMCACHED = "memcached"
)

func NewCacheStorageFactory(libcache, server, database string) (interface{}, error) {
	switch libcache {
	case REDIS:
		return newCacheRedis(server, database)
	case MEMCACHED:
		return newCacheMemcached(server)
	default:
		return nil, fmt.Errorf("unrecognized cache storage [%s]", libcache)
	}
}