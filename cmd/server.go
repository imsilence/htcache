package cmd

import (
	"fmt"
	"htcache/cache"
	"htcache/server"
	"github.com/spf13/cobra"
	_ "htcache/cache/provider"
	_ "htcache/server/provider"
)

var ctype string
var stype string
var sport int

var serverCmd *cobra.Command = &cobra.Command{
	Use:   "server",
	Short: "htcache server",
	Long:  "htcache server",
	RunE: func(cmd *cobra.Command, args []string) error {
		server.New(stype, cache.New(ctype)).Listen(fmt.Sprintf(":%d", sport))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&ctype, "type", "t", "memory", "cache type [memory/rocksdb]")
	serverCmd.Flags().StringVarP(&stype, "server", "s", "http", "server type [http/tcp]")
	serverCmd.Flags().IntVarP(&sport, "port", "p", 8888, "server port")
}
