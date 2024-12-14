package cmd

import (
	testApp "github.com/AEnjoy/IoT-lubricant/internal/test"
	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
	"github.com/AEnjoy/IoT-lubricant/pkg/test"
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
			panic(testApp.Service(&testApp.Mini{}).App(cli, abort, init))
		},
		SilenceUsage:  false,
		SilenceErrors: false,
	}
	testCommand.Flags().Bool("has-inited", false, "Whether the agent has been initialized")
	return testCommand
}
