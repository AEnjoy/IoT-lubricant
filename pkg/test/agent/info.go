package agent

import (
	"context"
	"fmt"

	testMeta "github.com/AEnjoy/IoT-lubricant/pkg/test"
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
)

func TestGetAgentInfo(cli agent.EdgeServiceClient) *testMeta.Result {
	fmt.Println("Test_GetAgentInfo:")
	_, err := cli.SetAgent(context.Background(), &agent.SetAgentRequest{
		AgentID: testMeta.AgentID,
		AgentInfo: &agent.AgentInfo{
			AgentID:     testMeta.AgentID,
			Description: &testMeta.AgentID,
		},
	})
	if err != nil {
		return &testMeta.Result{Success: false, Message: err.Error()}
	}
	fmt.Print("--Test set agent info: ")
	info, err := cli.GetAgentInfo(context.Background(), &agent.GetAgentInfoRequest{
		AgentID: testMeta.AgentID,
	})
	if err != nil {
		return &testMeta.Result{Success: false, Message: err.Error()}
	}
	if info.GetAgentInfo().GetDescription() != testMeta.AgentID {
		fmt.Println("failed")
		return &testMeta.Result{Success: false, Message: "get agent info failed"}
	}
	fmt.Println("Success")
	return &testMeta.Result{Success: true}
}
