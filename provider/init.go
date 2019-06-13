package provider

import (
	"htcache/cache"
	"htcache/provider/memory"
	"htcache/provider/rocksdb"
)

func init() {
	cache.Register("memory", memory.New)
	cache.Register("rocksdb", rocksdb.New)
}