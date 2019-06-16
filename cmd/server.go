package cmd

import (
	"fmt"
	"htcache/cache"
	_ "htcache/cache/provider"
	"htcache/server"
	_ "htcache/server/provider"
	"htcache/server/cluster"
	_ "htcache/server/cluster/provider"
	"github.com/spf13/cobra"
)

var (
	ctype string
	stype string
	sport int
	scluster string
	saddr string
)

var serverCmd *cobra.Command = &cobra.Command{
	Use:   "server",
	Short: "htcache server",
	Long:  "htcache server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cache, err := cache.NewCache(ctype)
		if err != nil {
			return err
		}
		node, err := cluster.NewNode("gossip", saddr, scluster)
		if err != nil {
			return err
		}
		server, err := server.NewServer(stype, node, cache)
		if err != nil {
			return err
		}
		return server.Listen(fmt.Sprintf("%s:%d", saddr, sport))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&scluster, "cluster", "c", "", "cluster addr")
	serverCmd.Flags().StringVarP(&saddr, "addr", "l", "127.0.0.1", "listen addr")
	serverCmd.Flags().IntVarP(&sport, "port", "p", 8888, "server port")
	serverCmd.Flags().StringVarP(&ctype, "type", "t", "memory", "cache type [memory/rocksdb]")
	serverCmd.Flags().StringVarP(&stype, "server", "s", "http", "server type [http/tcp]")

}
