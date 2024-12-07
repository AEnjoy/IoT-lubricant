package main

import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/cmd/test/agent/cmd"
	"github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "agent-test",
		Short: "client is an agent-test-client",
		Long:  `client is an agent-test-client`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Help())
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}
func init() {
	rootCmd.PersistentFlags().String("agent-address", fmt.Sprintf("%s:%d", "127.0.0.1", _default.AgentGrpcPort), "agent-grpc-address")
	rootCmd.PersistentFlags().String("agent-id", "", "test agent id")
	rootCmd.PersistentFlags().Bool("auto-abort", true, "When failed, automatically terminate subsequent tests.")
	rootCmd.AddCommand(cmd.MiniTest())
}
