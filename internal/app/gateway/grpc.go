package gateway

import (
	"io"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	taskTypes "github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/AEnjoy/IoT-lubricant/protobuf/meta"
	json "github.com/bytedance/sonic"
)

var _ = &meta.Ping{
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
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		var c types.TaskCommand
		var taskId string

		switch task := resp.GetTask().(type) {
		case *core.Task_GatewayGetTaskResponse:
			content := task.GatewayGetTaskResponse.GetMessage().GetContent()
			taskId = task.GatewayGetTaskResponse.GetMessage().GetTaskId()
			err = json.Unmarshal(content, &c)
			if err != nil {
				return err
			}
		case *core.Task_CorePushTaskRequest:
			content := task.CorePushTaskRequest.GetMessage().GetContent()
			taskId = task.CorePushTaskRequest.GetMessage().GetTaskId()
			err = json.Unmarshal(content, &c)
			if err != nil {
				return err
			}
		}

		switch c.ID {
		case taskTypes.OperationAddAgent:
			var req model.CreateAgentRequest
			err := json.Unmarshal(c.Data, &req)
			if err != nil {
				return err
			}

			_, err = HandelCreateAgentRequest(&req)
			if err != nil {
				return err
			}
			//result, _ := json.Marshal(response)

			resp := &core.Task{
				ID: taskId,
				Task: &core.Task_CorePushTaskResponse{
					CorePushTaskResponse: &core.CorePushTaskResponse{
						//Message: &core.TaskDetail{
						//	Content: result,
						//	TaskId:  taskId,
						//},
					},
				},
			}
			_ = task.Send(resp)
		case taskTypes.OperationRemoveAgent:
			a.agentRemove("reserve a seat")
			panic("not implemented")
		case taskTypes.OperationNil:

		default:
			panic("unhandled default case")
		}
	}

}
