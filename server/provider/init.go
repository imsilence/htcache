package provider

import (
	"htcache/server"
	"htcache/server/provider/http"
	"htcache/server/provider/tcp"
)

func init() {
	server.Register("http", http.New)
	server.Register("tcp", tcp.New)
}