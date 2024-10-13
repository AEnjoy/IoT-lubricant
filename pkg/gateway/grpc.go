package gateway

import (
	"encoding/json"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
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
		var t model.Command
		err = json.Unmarshal(resp.GetContent(), &t)
		if err != nil {
			return err
		}
		switch t.ID {
		case model.Command_RemoveAgent:
			a.removeAgent(t.Data)
		case model.Command_nil:

		default:
			panic("unhandled default case")
		}
	}

}
