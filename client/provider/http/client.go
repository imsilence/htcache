package http

import (
	"bytes"
	"fmt"
	"htcache/client"
	"io/ioutil"
	"net/http"
	"errors"
)

type Client struct {
	*http.Client
	Addr string
}

func New(addr string) (client.Client, error) {
	return &Client{
		Client: &http.Client{
			Transport: &http.Transport{},
		},
		Addr: addr,
	}, nil
}

func (c *Client) Run(command *client.Command) error {
	url := fmt.Sprintf("http://%s/cache/%s/", c.Addr, command.Key)

	switch command.Name {
	case "get":
		response, err := c.Get(url)
		if err != nil {
			command.Error = err
		} else if response.StatusCode == http.StatusNotFound {
			command.Value = []byte{}
		} else if response.StatusCode != http.StatusOK {
			command.Error = fmt.Errorf("request error: %d", response.StatusCode)
		} else {
			if bytes, err := ioutil.ReadAll(response.Body); err == nil {
				command.Value = bytes
			} else {
				command.Error = err
			}
		}
	case "set":
		response, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(command.Value))
		if err != nil {
			command.Error = err
		} else if response.StatusCode != http.StatusOK {
			command.Error = fmt.Errorf("request error: %d", response.StatusCode)
		}
	case "delete":
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			command.Error = err
		} else {
			response, err := c.Do(request)
			if err != nil {
				command.Error = err
			} else if response.StatusCode != http.StatusOK {
				command.Error = fmt.Errorf("request error: %d", response.StatusCode)
			}
		}
	}
	return nil
}

func (c *Client) Pipeline(commands []*client.Command) error {
	return errors.New("http pipeline run not implement")
}

func (c *Client) Close() error {
	return nil
}