package memory

import (
	"htcache/cache"
	"sync"
	"time"
)

type Value struct {
	value   []byte
	created time.Time
}

type Cache struct {
	cache map[string]Value
	mutex sync.RWMutex
	cache.Stat
	ttl time.Duration
}

func New(ttl time.Duration) (cache.Cache, error) {
	cache := &Cache{
		cache: make(map[string]Value),
		mutex: sync.RWMutex{},
		Stat:  cache.Stat{},
		ttl:   ttl,
	}
	if int(ttl) > 0 {
		go func() {
			for now := range time.Tick(time.Second) {
				cache.mutex.RLock()
				for k, v := range cache.cache {
					if v.created.Add(ttl).Before(now) {
						cache.mutex.RUnlock()
						cache.Del(k)
						cache.mutex.RLock()
					}
				}
				cache.mutex.RUnlock()
			}
		}()
	}
	return cache, nil
}

func (c *Cache) Get(key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if v, ok := c.cache[key]; ok {
		return v.value, nil
	} else {
		return nil, cache.NotFound
	}
}

func (c *Cache) Set(key string, value []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.cache[key]; ok {
		c.Stat.Del(key, v.value)
	}
	c.cache[key] = Value{value, time.Now()}
	c.Stat.Add(key, value)
	return nil
}

func (c *Cache) Del(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.cache[key]; ok {
		delete(c.cache, key)
		c.Stat.Del(key, v.value)
	}
	return nil
}

func (c *Cache) GetStat() cache.Stat {
	return c.Stat
}
