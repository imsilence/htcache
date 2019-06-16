package cmd

import (
	"fmt"
	"htcache/cache"
	_ "htcache/cache/provider"
	"htcache/server"
	"htcache/server/cluster"
	_ "htcache/server/cluster/provider"
	_ "htcache/server/provider"

	"github.com/spf13/cobra"
)

var (
	ctype    string
	stype    string
	sport    int
	scluster string
	scaddr   string
	scport   int
	saddr    string
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
		node, err := cluster.NewNode("gossip", scluster, fmt.Sprintf("%s:%d", saddr, sport), fmt.Sprintf("%s:%d", scaddr, scport))
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
	serverCmd.Flags().StringVarP(&scaddr, "caddr", "H", "127.0.0.1", "cluster manager addr")
	serverCmd.Flags().IntVarP(&scport, "cport", "P", 8889, "cluster manager port")
	serverCmd.Flags().StringVarP(&ctype, "type", "t", "memory", "cache type [memory/rocksdb]")
	serverCmd.Flags().StringVarP(&stype, "server", "s", "http", "server type [http/tcp]")

}
