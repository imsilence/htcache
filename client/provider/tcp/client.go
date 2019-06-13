package tcp

import (
	"htcache/client"
)

type Client struct {
	addr string
}


func New(addr string) client.Client {
	return &Client{addr}
}


func (c *Client)Run(command *Command) {

}

func (c *Client)Pipeline(commands []*Command) {

}