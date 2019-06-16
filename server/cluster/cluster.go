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

func NewNode(name, cluster, addr, maddr string) (Node, error) {
	new, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("node %s is unregister", name)
	}

	node, err := new(addr, cluster)
	if err != nil {
		return node, err
	}

	server, err := NewServer(node)
	if err != nil {
		return node, err
	}
	go func() {
		server.Listen(maddr)
	} ()
	return node, nil

}
