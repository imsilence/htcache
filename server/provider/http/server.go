package http

import (
	"htcache/server/cluster"
	"htcache/cache"
	"htcache/server"
	"log"
	"net/http"
)

type Server struct {
	cluster.Node
	cache.Cache
}

func New(n cluster.Node, c cache.Cache) (server.Server, error) {
	return &Server{n, c}, nil
}

func (s *Server) Listen(addr string) error {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status/", s.statusHandler())

	log.Printf("Started HTTP Server, Listen On: %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) cacheHandler() http.Handler {
	return &CacheHandler{s}
}

func (s *Server) statusHandler() http.Handler {
	return &StatusHandler{s}
}
