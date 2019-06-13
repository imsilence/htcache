package client


type Command struct {
	Name string
	Key string
	Value []byte
	Error error
}

type Client interface {
	Run(*Command)
	Pipline([]*Command)
}