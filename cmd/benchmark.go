package cmd

import (
	"fmt"
	"htcache/benchmark"
	"os"

	"github.com/spf13/cobra"
)

var (
	bprotocol   string
	bhost       string
	bport       int
	bconcurrent int
	btotal      int
	boperation  string
)

var benchmarkCmd *cobra.Command = &cobra.Command{
	Use:   "benchmark",
	Short: "htcache benchmark tools",
	Long:  "htcache benchmark tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		bm := benchmark.NewBenchmark(
			bprotocol,
			fmt.Sprintf("%s:%d", bhost, bport),
			bconcurrent,
			btotal,
			boperation,
		)
		bm.Execute()
		bm.Output(os.Stdin)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(benchmarkCmd)

	benchmarkCmd.Flags().StringVarP(&bprotocol, "proto", "p", "http", "htcache protocol type [http/tcp]")
	benchmarkCmd.Flags().StringVarP(&bhost, "host", "H", "localhost", "htcache server addr")
	benchmarkCmd.Flags().IntVarP(&bport, "port", "P", 8888, "htcache server port")
	benchmarkCmd.Flags().IntVarP(&bconcurrent, "concurrent", "c", 10, "concurrent goroutine number")
	benchmarkCmd.Flags().IntVarP(&btotal, "total", "t", 1000, "operation total count")
	benchmarkCmd.Flags().StringVarP(&boperation, "operation", "o", "mixed", "operation type[set/get/mixed]")
}
