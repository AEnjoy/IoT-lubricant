package gateway

import (
	"encoding/json"

	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	taskTypes "github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
)

var _ = &core.Ping{
	Flag: 0,
}

func (a *app) grpcApp() error {
	// todo: not all implemented yet
	task, err := a.grpcClient.GetTask(a.ctrl)
	if err != nil {
		return err
	}
	for {
		resp, err := task.Recv()
		if err != nil {
			return err
		}

		var c types.TaskCommand
		switch task := resp.GetTask().(type) {
		case *core.Task_GatewayGetTaskResponse:
			content := task.GatewayGetTaskResponse.GetMessage().GetContent()
			err = json.Unmarshal(content, &c)
			if err != nil {
				return err
			}
		case *core.Task_CorePushTaskRequest:
			content := task.CorePushTaskRequest.GetMessage().GetContent()
			err = json.Unmarshal(content, &c)
			if err != nil {
				return err
			}
		}

		switch c.ID {
		case taskTypes.OperationRemoveAgent:
			//a.removeAgent(c.Data)
		case taskTypes.OperationNil:

		default:
			panic("unhandled default case")
		}
	}

}
