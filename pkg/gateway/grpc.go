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
		var t types.Command
		task := resp.GetTask()
		switch t := task.(type) {
		case *core.Task_GatewayGetTaskResponse:
			content := t.GatewayGetTaskResponse.GetMessage().GetContent()
			err = json.Unmarshal(content, &t)
			if err != nil {
				return err
			}
		case *core.Task_CorePushDataRequest:
			content := t.CorePushDataRequest.GetMessage().GetContent()
			err = json.Unmarshal(content, &t)
			if err != nil {
				return err
			}
		}
		switch t.ID {
		case types.Command_RemoveAgent:
			a.removeAgent(t.Data)
		case types.Command_nil:

		default:
			panic("unhandled default case")
		}
	}

}
