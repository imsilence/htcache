package cache

import (
	"errors"
	"fmt"
)

var NotFound error = errors.New("Not Found")

type NewFunc func() Cache

var providers map[string]NewFunc = make(map[string]NewFunc)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Del(key string) error
	GetStat() Stat
}

func Register(name string, new NewFunc) {
	if _, ok := providers[name]; ok {
		panic(fmt.Sprintf("cache %s is registered", name))
	}
	providers[name] = new
}

func New(name string) Cache {
	if new, ok := providers[name]; ok {
		return new()
	}
	panic(fmt.Sprintf("cache %s is not unregister", name))
	return nil
}
