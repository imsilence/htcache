package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clientCmd *cobra.Command = &cobra.Command{
	Use:   "client",
	Short: "htcache client",
	Long:  "htcache client",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("c")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

}
