package main

import (
	"htcache/http"
	"htcache/provider/memory"
	"htcache/tcp"
)

func main() {
	cache := memory.New()
	stop := make(chan int, 0)
	go func() {
		http.New(cache).Listen(":8888")
		stop <- 1
	}()
	go func() {
		tcp.New(cache).Listen(":8889")
		stop <- 2
	}()
	<-stop
}
