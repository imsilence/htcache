package http

import (
	"encoding/json"
	"htcache/cache"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type CacheHandler struct {
	*Server
}

func (handler *CacheHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	paths := strings.Split(request.URL.EscapedPath(), "/")
	if len(paths) < 3 || len(strings.TrimSpace(paths[2])) == 0 {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	key := strings.TrimSpace(paths[2])

	switch request.Method {
	case http.MethodGet:
		if bytes, err := handler.Get(key); err == nil {
			log.Printf("Get Key: %s, Value: %v", key, bytes)
			response.Write(bytes)
		} else if err == cache.NotFound {
			log.Printf("Get Key: %s, Value: Empty", key)
			response.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("Get Key: %s, Error: %v", key, err)
			response.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPost:
		bytes, err := ioutil.ReadAll(request.Body)
		if err == nil {
			if err := handler.Set(key, bytes); err == nil {
				log.Printf("Set Key: %s, Value: %v", key, bytes)
				response.WriteHeader(http.StatusOK)
			} else {
				log.Printf("Set Key: %s, Value: %v, Error: %v", key, bytes, err)
				response.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			response.WriteHeader(http.StatusBadRequest)
		}
	case http.MethodDelete:
		if err := handler.Del(key); err == nil {
			log.Printf("Delete Key: %s", key)
			response.WriteHeader(http.StatusOK)
		} else {
			log.Printf("Delete Key: %s, Error: %v", key, err)
			response.WriteHeader(http.StatusInternalServerError)
		}
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type StatusHandler struct {
	*Server
}

func (handler *StatusHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
	if bytes, err := json.Marshal(handler.GetStat()); err == nil {
		response.Write(bytes)
	} else {
		response.WriteHeader(http.StatusInternalServerError)
	}
}
