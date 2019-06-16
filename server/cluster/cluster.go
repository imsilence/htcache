package cluster

import "fmt"

type Node interface {
	IsProcess(key string) (string, bool)
	Members() []string
	Addr() string
}

type NewFunc func(addr, cluster string) (Node, error)

var providers map[string]NewFunc = make(map[string]NewFunc)

func Register(name string, new NewFunc) {
	if _, ok := providers[name]; ok {
		panic(fmt.Sprintf("node %s is registered", name))
	}
	providers[name] = new
}

func NewNode(name, addr, cluster string) (Node, error) {
	if new, ok := providers[name]; ok {
		return new(addr, cluster)
	}
	return nil, fmt.Errorf("node %s is unregister", name)
}
