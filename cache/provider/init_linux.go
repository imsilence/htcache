package cache

import (
	"htcache/cache"
	"htcache/cache/provider/memory"
	"htcache/cache/provider/rocksdb"
)

func init() {
	cache.Register("memory", memory.New)
	cache.Register("rocksdb", rocksdb.New)
}