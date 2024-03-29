package cache

import (
	"errors"
	"fmt"
	"time"
)

var NotFound error = errors.New("Not Found")

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Del(key string) error
	GetStat() Stat
}

type NewFunc func(time.Duration) (Cache, error)

var providers map[string]NewFunc = make(map[string]NewFunc)

func Register(name string, new NewFunc) {
	if _, ok := providers[name]; ok {
		panic(fmt.Sprintf("cache %s is registered", name))
	}
	providers[name] = new
}

func NewCache(name string, ttl time.Duration) (Cache, error) {
	if new, ok := providers[name]; ok {
		return new(ttl)
	}
	return nil, fmt.Errorf("cache %s is not unregister", name)
}
