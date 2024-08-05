package memoryCache

import "sync"

type CacheDriver struct {
	Cache sync.Map
}

func GetCacheDriver() (*CacheDriver, error) {
	return &CacheDriver{Cache: sync.Map{}}, nil
}
