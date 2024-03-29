package client

import (
	"fmt"
)

type NewFunc func(string) (Client, error)

var providers map[string]NewFunc = make(map[string]NewFunc)

type Command struct {
	Name  string
	Key   string
	Value []byte
	Error error
}

func NewCommand(name, key string, value []byte) *Command {
	return &Command{
		Name:  name,
		Key:   key,
		Value: value,
		Error: nil,
	}
}

func (c *Command) String() string {
	return fmt.Sprintf("Name: %s, Key: %s, Value: %v, Error: %v", c.Name, c.Key, c.Value, c.Error)
}

type Client interface {
	Run(*Command) error
	Pipeline([]*Command) error
	Close() error
}

func Register(name string, new NewFunc) {
	if _, ok := providers[name]; ok {
		panic(fmt.Sprintf("client %s is registered", name))
	}
	providers[name] = new
}

func NewClient(name string, addr string) (Client, error) {
	if new, ok := providers[name]; ok {
		return new(addr)
	}
	return nil, fmt.Errorf("client %s is unregister", name)
}
