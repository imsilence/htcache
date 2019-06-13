package client

import (
	"fmt"
)

type NewFunc func() Client

var providers map[string]NewFunc = make(map[string]NewFunc)

type Command struct {
	Name string
	Key string
	Value []byte
	Error error
}

type Client interface {
	Run(*Command)
	Pipeline([]*Command)
}

func Register(name string, new NewFunc) {
	if _, ok := providers[name]; ok {
		panic(fmt.Sprintf("client %s is registered", name))
	}
	providers[name] = new
}

func New(name string) Client {
	if new, ok:= providers[name]; ok {
		return new()
	}
	panic(fmt.Sprintf("client %s is unregister", name))
	return nil
}