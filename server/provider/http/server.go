package http

import (
	"htcache/cache"
	"net/http"
	"htcache/server"
)

type Server struct {
	cache.Cache
}

func New(c cache.Cache) server.Server {
	return &Server{c}
}

func (s *Server) Listen(addr string) {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status/", s.statusHandler())
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

func (s *Server) cacheHandler() http.Handler {
	return &CacheHandler{s}
}

func (s *Server) statusHandler() http.Handler {
	return &StatusHandler{s}
}
