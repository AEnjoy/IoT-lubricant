package agent

import (
	"fmt"

	testMeta "github.com/AEnjoy/IoT-lubricant/pkg/test"
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
)

func TestStartGather(cli agent.EdgeServiceClient) *testMeta.Result {
	fmt.Println("Test_StartGather:")
	// todo
	// 没有设置doc->Invalid internal configuration
	// 已经启动->Gather is working now
	panic("implement me")
}
func TestStopGather(cli agent.EdgeServiceClient) *testMeta.Result {
	fmt.Println("Test_StopGather:")
	// todo
	// 已经启动->success
	// 没有启动->Gather is not working now
	panic("implement me")
}
