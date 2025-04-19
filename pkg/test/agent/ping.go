package agent

import (
	"context"
	"fmt"

	testMeta "github.com/aenjoy/iot-lubricant/pkg/test"
	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
	agentMeta "github.com/aenjoy/iot-lubricant/protobuf/meta"
)

func TestPing(cli agentpb.EdgeServiceClient) *testMeta.Result {
	fmt.Println("Test Ping:")
	pingResult, err := cli.Ping(context.Background(), &agentMeta.Ping{})
	if err != nil {
		fmt.Println("Test Ping failed:")
		return &testMeta.Result{Success: false, Message: err.Error()}
	}
	if pingResult.GetFlag() != 2 {
		return &testMeta.Result{Success: false, Message: "error ping return value"}
	}
	return &testMeta.Result{Success: true}
}
