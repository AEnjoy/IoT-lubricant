package agent

import (
	"context"
	"fmt"
	"net/http"

	testMeta "github.com/aenjoy/iot-lubricant/pkg/test"
	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
)

func TestRegisterGateway(cli agentpb.EdgeServiceClient) *testMeta.Result {
	fmt.Println("Test_RegisterGateway:")
	fmt.Print("--Test error target:")
	registerGatewayResponse, err := cli.RegisterGateway(context.Background(),
		&agentpb.RegisterGatewayRequest{},
	)
	if err != nil {
		return &testMeta.Result{Success: false, Message: err.Error()}
	}
	if errIsTargetNotEqual(registerGatewayResponse.GetInfo()) {
		fmt.Println("Success")
	}
	fmt.Print("--Test success register:")
	registerGatewayResponse, err = cli.RegisterGateway(context.Background(),
		&agentpb.RegisterGatewayRequest{
			AgentID:   testMeta.AgentID,
			GatewayID: testMeta.GatewayID,
		},
	)
	if err != nil {
		return &testMeta.Result{Success: false, Message: err.Error()}
	}
	if registerGatewayResponse.GetInfo().Code != http.StatusOK {
		fmt.Println("failed")
		return &testMeta.Result{Success: false, Message: "register gateway failed"}
	}
	fmt.Println("Success")
	return &testMeta.Result{Success: true}
}
