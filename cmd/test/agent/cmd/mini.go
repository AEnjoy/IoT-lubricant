package cmd

import (
	"github.com/aenjoy/iot-lubricant/pkg/edge"
	"github.com/aenjoy/iot-lubricant/pkg/test"
	test2 "github.com/aenjoy/iot-lubricant/services/test"
	"github.com/spf13/cobra"
)

func MiniTest() *cobra.Command {
	var testCommand = &cobra.Command{
		Use:   "mini",
		Short: "Automatically execute the smallest tests",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("agent-id")
			test.AgentID = id
			address, _ := cmd.Flags().GetString("agent-address")
			if address == "" {
				panic("Please specify the agent address")
			}
			cli, err := edge.NewAgentCli(address)
			if err != nil {
				panic(err.Error())
			}
			abort, _ := cmd.Flags().GetBool("auto-abort")
			init, _ := cmd.Flags().GetBool("has-inited")
			panic(test2.Service(&test2.Mini{}).App(cli, abort, init))
		},
		SilenceUsage:  false,
		SilenceErrors: false,
	}
	testCommand.Flags().Bool("has-inited", false, "Whether the agent has been initialized")
	return testCommand
}
