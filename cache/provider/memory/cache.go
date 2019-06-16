package memory

import (
	"htcache/cache"
	"sync"
)

type Cache struct {
	cache map[string][]byte
	mutex sync.RWMutex
	cache.Stat
}

func New() (cache.Cache, error) {
	return &Cache{
		cache: make(map[string][]byte),
		mutex: sync.RWMutex{},
		Stat:  cache.Stat{},
	}, nil
}

func (c *Cache) Get(key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if v, ok := c.cache[key]; ok {
		return v, nil
	} else {
		return nil, cache.NotFound
	}
}

func (c *Cache) Set(key string, value []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.cache[key]; ok {
		c.Stat.Del(key, v)
	}
	c.cache[key] = value
	c.Stat.Add(key, value)
	return nil
}

func (c *Cache) Del(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.cache[key]; ok {
		delete(c.cache, key)
		c.Stat.Del(key, v)
	}
	return nil
}

func (c *Cache) GetStat() cache.Stat {
	return c.Stat
}
