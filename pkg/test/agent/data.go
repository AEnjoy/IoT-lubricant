package agent

import (
	"fmt"

	testMeta "github.com/AEnjoy/IoT-lubricant/pkg/test"
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
)

func TestGetDataStream(cli agent.EdgeServiceClient) *testMeta.Result {
	// todo implement me
	panic("implement me")
}
func TestGetGatherData(cli agent.EdgeServiceClient, testType int) *testMeta.Result {
	fmt.Println("Test_GetGatherData:")
	//
	switch testType {
	// todo
	// 0:未配置文档/api 采集失败
	// 1:配置文档/api,未启动采集 采集为空 显示data is not ready
	// 2:配置文档/api,已启动采集 采集成功
	}
	return &testMeta.Result{Success: false, Message: "not support test type"}
}

func TestSendHttpRequest(cli agent.EdgeServiceClient, testType int) *testMeta.Result {
	fmt.Println("Test_SendHttpRequest:")
	switch testType {
	// todo
	// 0:未配置文档/api Invalid internal configuration
	// 1:配置文档/api 无效路径 Invalid path
	// 2:配置文档/api 错误入参: 抛出err
	// 3:配置文档/api 正确入参: 获取返回值
	}
	return &testMeta.Result{Success: false, Message: "not support test type"}
}
