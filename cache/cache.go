package cache

import (
	"errors"
	"fmt"
)

var NotFound error = errors.New("Not Found")

var providers map[string]func() Cache = make(map[string]func() Cache)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Del(key string) error
	GetStat() Stat
}

func Register(ctype string, new func() Cache) {
	if _, ok := providers[ctype]; ok {
		panic(fmt.Sprintf("cache %s is exists", ctype))
	}
	providers[ctype] = new
}

func New(ctype string) Cache {
	if new, ok := providers[ctype]; ok {
		return new()
	}
	panic(fmt.Sprintf("cache %s is not defined", ctype))
	return nil
}
