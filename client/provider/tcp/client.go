package tcp

import (
	"htcache/client"
	"htcache/cache"
	"net"
	"fmt"
	"bufio"
	"strconv"
	"io"
	"strings"
	"errors"
)

type Client struct {
	net.Conn
	Addr string
}

func New(addr string) (client.Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn, addr}, nil
}

func (c *Client) Run(command *client.Command) error {
	reader := bufio.NewReader(c)
	switch command.Name {
	case "get":
		fmt.Fprintf(c, "G%d %s", len(command.Key), command.Key)
		value, err := c.readResult(reader)
		if err == nil {
			command.Value = value
		} else if err.Error() == cache.NotFound.Error() {
			command.Value = nil
		} else {
			command.Error = err
		}
	case "set":
		fmt.Fprintf(c, "S%d %d %s%s", len(command.Key), len(command.Value), command.Key, command.Value)
		command.Value, command.Error = c.readResult(reader)
	case "delete", "del":
		fmt.Fprintf(c, "D%d %s", len(command.Key), command.Key)
		command.Value, command.Error = c.readResult(reader)
	}
	return nil
}

func (s *Client) readLen(reader *bufio.Reader) (int, error) {
	cxt, err := reader.ReadString(' ')
	if err != nil {
		return 0, err
	}
	len, err := strconv.Atoi(strings.TrimSpace(cxt))
	if err != nil {
		return 0, err
	}
	return len, nil
}

func (s *Client) readResult(reader *bufio.Reader) ([]byte, error) {
	len, err := s.readLen(reader)
	if err != nil {
		return nil, err
	}
	rlen := len
	if rlen < 0 {
		rlen = -rlen
	}
	bytes := make([]byte, rlen)
	_, err = io.ReadFull(reader, bytes)
	if err != nil {
		return nil, err
	}
	if len > 0 {
		return bytes, nil
	} else {
		return nil, errors.New(string(bytes))
	}

}

func (c *Client) Pipeline(commands []*client.Command) error {
	reader := bufio.NewReader(c)
	for _, command:= range commands {
		switch command.Name {
		case "get":
			fmt.Fprintf(c, "G%d %s", len(command.Key), command.Key)
		case "set":
			fmt.Fprintf(c, "S%d %d %s%s", len(command.Key), len(command.Value), command.Key, command.Value)
		case "delete", "del":
			fmt.Fprintf(c, "D%d %s", len(command.Key), command.Key)
		}
	}
	for _, command := range commands {
		switch command.Name {
		case "get":
			value, err := c.readResult(reader)
			if err == nil {
				command.Value = value
			} else if err.Error() == cache.NotFound.Error() {
				command.Value = nil
			} else {
				command.Error = err
			}
		case "set":
			command.Value, command.Error = c.readResult(reader)
		case "delete", "del":
			command.Value, command.Error = c.readResult(reader)
		}
	}
	return nil
}
