package server

import (
	"fmt"
	"htcache/cache"
	"htcache/server/cluster"
)

type Server interface {
	Listen(addr string) error
}

type NewFunc func(cluster.Node, cache.Cache) (Server, error)

var providers map[string]NewFunc = make(map[string]NewFunc)

func Register(name string, new NewFunc) {
	if _, ok := providers[name]; ok {
		panic(fmt.Sprintf("server %s is registered", name))
	}
	providers[name] = new
}

func NewServer(name string, node cluster.Node, cache cache.Cache) (Server, error) {
	if new, ok := providers[name]; ok {
		return new(node, cache)
	}
	return nil, fmt.Errorf("server %s is unregister", name)
}
