package gateway

import (
	"encoding/json"

	"github.com/AEnjoy/IoT-lubricant/pkg/types"
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

		var c types.Command
		switch task := resp.GetTask().(type) {
		case *core.Task_GatewayGetTaskResponse:
			content := task.GatewayGetTaskResponse.GetMessage().GetContent()
			err = json.Unmarshal(content, &c)
			if err != nil {
				return err
			}
		case *core.Task_CorePushDataRequest:
			content := task.CorePushDataRequest.GetMessage().GetContent()
			err = json.Unmarshal(content, &c)
			if err != nil {
				return err
			}
		}

		switch c.ID {
		case types.Command_RemoveAgent:
			a.removeAgent(c.Data)
		case types.Command_nil:

		default:
			panic("unhandled default case")
		}
	}

}
