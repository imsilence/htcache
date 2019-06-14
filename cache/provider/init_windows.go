package cache

import (
	"htcache/cache"
	"htcache/cache/provider/memory"
)

func init() {
	cache.Register("memory", memory.New)
}