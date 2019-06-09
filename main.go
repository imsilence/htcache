package main

import (
	"htcache/http"
	"htcache/provider/memory"
)

func main() {
	http.New(memory.New()).Listen(":8888")
}
