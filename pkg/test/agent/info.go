package agent

import (
	"context"
	"fmt"

	testMeta "github.com/aenjoy/iot-lubricant/pkg/test"
	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
)

func TestGetAgentInfo(cli agentpb.EdgeServiceClient) *testMeta.Result {
	fmt.Println("Test_GetAgentInfo:")
	_, err := cli.SetAgent(context.Background(), &agentpb.SetAgentRequest{
		AgentID: testMeta.AgentID,
		AgentInfo: &agentpb.AgentInfo{
			AgentID:     testMeta.AgentID,
			Description: &testMeta.AgentID,
		},
	})
	if err != nil {
		return &testMeta.Result{Success: false, Message: err.Error()}
	}
	fmt.Print("--Test set agent info: ")
	info, err := cli.GetAgentInfo(context.Background(), &agentpb.GetAgentInfoRequest{
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
