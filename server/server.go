package server

import (
	"fmt"
	"htcache/cache"
)

type NewFunc func(cache.Cache) Server

var providers map[string]NewFunc = make(map[string]NewFunc)

type Server interface {
	Listen(addr string)
}

func Register(name string, new NewFunc) {
	if _, ok := providers[name]; ok {
		panic(fmt.Sprintf("server %s is registered", name))
	}
	providers[name] = new
}

func NewServer(name string, cache cache.Cache) Server {
	if new, ok := providers[name]; ok {
		return new(cache)
	}
	panic(fmt.Sprintf("server %s is unregister", name))
	return nil
}
