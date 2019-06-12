package main

import (
	"flag"
	"fmt"
	"htcache/cache"
	"htcache/http"
	"htcache/provider/memory"
	"htcache/provider/rocksdb"
	"htcache/tcp"
	"os"
)

func main() {
	ctype := flag.String("type", "memory", "cache type [memory/rocksdb]")
	phttp := flag.Int("http", 8888, "http server port")
	ptcp := flag.Int("tcp", 8889, "tcp server port")
	help := flag.Bool("help", false, "help")
	h := flag.Bool("h", false, "help")

	flag.Usage = func() {
		fmt.Println("Usage: htcache [-type memory] [-http 8888] [-tcp 8889]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *h || *help {
		flag.Usage()
		os.Exit(0)
	}

	var cache cache.Cache
	switch *ctype {
	case "rocksdb":
		cache = rocksdb.New()
	default:
		cache = memory.New()
	}
	stop := make(chan int, 0)
	go func() {
		http.New(cache).Listen(fmt.Sprintf(":%d", *phttp))
		stop <- 1
	}()
	go func() {
		tcp.New(cache).Listen(fmt.Sprintf(":%d", *ptcp))
		stop <- 2
	}()
	<-stop
}
