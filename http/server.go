package http

import (
	"htcache/cache"
	"net/http"
)

type Server struct {
	cache.Cache
}

func New(c cache.Cache) *Server {
	return &Server{c}
}

func (s *Server) Listen(addr string) {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status/", s.statusHandler())
	http.ListenAndServe(addr, nil)
}

func (s *Server) cacheHandler() http.Handler {
	return &CacheHandler{s}
}

func (s *Server) statusHandler() http.Handler {
	return &StatusHandler{s}
}
