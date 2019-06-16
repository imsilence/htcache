package cmd

import (
	"fmt"
	_ "htcache/client/provider"

	"htcache/client"

	"github.com/spf13/cobra"
)

var (
	cprotocol  string
	chost      string
	cport      int
	coperation string
	ckey       string
	cvalue     string
)

var clientCmd *cobra.Command = &cobra.Command{
	Use:   "client",
	Short: "htcache client",
	Long:  "htcache client",
	RunE: func(cmd *cobra.Command, args []string) error {
		command := client.NewCommand(coperation, ckey, []byte(cvalue))
		cli, err := client.NewClient(cprotocol, fmt.Sprintf("%s:%d", chost, cport))
		if err != nil {
			return err
		}
		defer cli.Close()
		cli.Run(command)
		fmt.Println(command)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&cprotocol, "proto", "p", "http", "htcache protocol type [http/tcp]")
	clientCmd.Flags().StringVarP(&chost, "host", "H", "localhost", "htcache server addr")
	clientCmd.Flags().IntVarP(&cport, "port", "P", 8888, "htcache server port")
	clientCmd.Flags().StringVarP(&coperation, "operation", "o", "get", "operation type[set/get/delete]")
	clientCmd.Flags().StringVarP(&ckey, "key", "k", "default", "cache key")
	clientCmd.Flags().StringVarP(&cvalue, "value", "v", "", "cache value")

}
