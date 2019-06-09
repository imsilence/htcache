package cache

import "errors"

var NotFound error = errors.New("Not Found")

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Del(key string) error
	GetStat() Stat
}
