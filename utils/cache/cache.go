package cache

import (
	"ess/utils/logging"
	"ess/utils/setting"
	"log"

	lru "github.com/hashicorp/golang-lru"
)

var cache *lru.ARCCache

func Setup() {
	var err error
	cache, err = lru.NewARC(setting.ServerSetting.CacheSize)
	if err != nil {
		log.Panicf("cache init fail :%v\n", err)
	}
}

// wipe the entire cache table
func ForceFlush() {
	logging.Info("force flush cache")
	cache.Purge()
}

// purge a key from the cache
func Remove(key string) {
	cache.Remove(key)
}

// accept a key and a function to return data (interface), if there isn't key-value in cache, it will call the
// function to get data and store it
func GetOrCreate(key string, f func() interface{}) interface{} {
	res, exist := cache.Get(key)

	if !exist {
		newValue := f()
		cache.Add(key, newValue)
		return newValue
	}
	return res
}

// put the key and value to the cache, and will force cover the old things
func Set(key string, value interface{}) {
	cache.Add(key, value)
}
