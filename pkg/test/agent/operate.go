package agent

import (
	"fmt"

	testMeta "github.com/aenjoy/iot-lubricant/pkg/test"
	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
)

func TestStartGather(cli agentpb.EdgeServiceClient) *testMeta.Result {
	fmt.Println("Test_StartGather:")
	// todo
	// 没有设置doc->Invalid internal configuration
	// 已经启动->Gather is working now
	panic("implement me")
}
func TestStopGather(cli agentpb.EdgeServiceClient) *testMeta.Result {
	fmt.Println("Test_StopGather:")
	// todo
	// 已经启动->success
	// 没有启动->Gather is not working now
	panic("implement me")
}
