package provider

import (
	"htcache/client"
	"htcache/client/provider/http"
	"htcache/client/provider/tcp"
)

func init() {
	client.Register("http", http.New)
	client.Register("tcp", tcp.New)
}