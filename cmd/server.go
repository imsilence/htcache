package cmd

import (
	"fmt"
	"htcache/cache"
	_ "htcache/cache/provider"
	"htcache/server"
	_ "htcache/server/provider"

	"github.com/spf13/cobra"
)

var (
	ctype string
	stype string
	sport int
)

var serverCmd *cobra.Command = &cobra.Command{
	Use:   "server",
	Short: "htcache server",
	Long:  "htcache server",
	RunE: func(cmd *cobra.Command, args []string) error {
		server.NewServer(stype, cache.NewCache(ctype)).Listen(fmt.Sprintf(":%d", sport))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&ctype, "type", "t", "memory", "cache type [memory/rocksdb]")
	serverCmd.Flags().StringVarP(&stype, "server", "s", "http", "server type [http/tcp]")
	serverCmd.Flags().IntVarP(&sport, "port", "p", 8888, "server port")
}
