package cmd

import (
	"fmt"
	"htcache/http"
	"htcache/cache"
	"htcache/tcp"
	"sync"
	"github.com/spf13/cobra"
	_ "htcache/provider"
)

var ctype string
var phttp int
var ptcp int

var serverCmd *cobra.Command = &cobra.Command{
	Use:   "server",
	Short: "htcache server",
	Long:  "htcache server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cache := cache.New(ctype)

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			http.New(cache).Listen(fmt.Sprintf(":%d", phttp))
			wg.Done()
		}()
		go func() {
			tcp.New(cache).Listen(fmt.Sprintf(":%d", ptcp))
			wg.Done()
		}()
		wg.Wait()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&ctype, "type", "t", "memory", "cache type [memory/rocksdb]")
	serverCmd.Flags().IntVarP(&phttp, "http", "H", 8888, "http server port")
	serverCmd.Flags().IntVarP(&ptcp, "tcp", "T", 8889, "tcp server port")
}
