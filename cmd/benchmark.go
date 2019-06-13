package cmd

import (
	"fmt"
	"htcache/benchmark"
	"os"

	"github.com/spf13/cobra"
)

var protocol string
var host string
var port int
var concurrent int
var total int
var operation string

var benchmarkCmd *cobra.Command = &cobra.Command{
	Use:   "benchmark",
	Short: "htcache benchmark tools",
	Long:  "htcache benchmark tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		bm := benchmark.NewBenchmark(
			protocol,
			fmt.Sprintf("%s:%d", host, port),
			concurrent,
			total,
			operation,
		)
		bm.Execute()
		bm.Output(os.Stdin)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(benchmarkCmd)

	benchmarkCmd.Flags().StringVarP(&protocol, "proto", "p", "http", "htcache protocol type [http/tcp]")
	benchmarkCmd.Flags().StringVarP(&host, "host", "H", "localhost", "htcache server addr")
	benchmarkCmd.Flags().IntVarP(&port, "port", "P", 8888, "htcache server port")
	benchmarkCmd.Flags().IntVarP(&concurrent, "concurrent", "c", 10, "concurrent goroutine number")
	benchmarkCmd.Flags().IntVarP(&total, "total", "t", 1000, "operation total count")
	benchmarkCmd.Flags().StringVarP(&operation, "operation", "o", "mixed", "operation type[set/get/mixed]")
}
