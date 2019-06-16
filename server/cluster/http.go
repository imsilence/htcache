package cluster

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Node
}

func NewServer(n Node) (*Server, error) {
	return &Server{n}, nil
}

func (s *Server) Listen(addr string) error {
	http.HandleFunc("/addr/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprint(response, s.Addr())
	})
	http.HandleFunc("/members/", func(response http.ResponseWriter, request *http.Request) {
		if bytes, err := json.Marshal(s.Members()); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Write(bytes)
		}

	})
	log.Printf("Started Cluster Manage HTTP Server, Listen On: %s", addr)
	return http.ListenAndServe(addr, nil)
}
