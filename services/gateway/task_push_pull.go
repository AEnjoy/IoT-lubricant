package gateway

import (
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"github.com/google/uuid"
)

func (a *app) _tasksAddToExecutor(taskDetail *corepb.TaskDetail, notice bool) *corepb.CorePushTaskResponse {
	var id string
	if taskDetail.TaskId == "" {
		id = uuid.NewString()
		taskDetail.TaskId = id
	}

	a.task.AddTask(taskDetail, notice)
	return &corepb.CorePushTaskResponse{
		TaskId: id,
	}
}
func (a *app) handelCorePushTaskAsync(task *corepb.Task_CorePushTaskRequest) *corepb.CorePushTaskResponse {
	return a._tasksAddToExecutor(task.CorePushTaskRequest.GetMessage(), true)
}

func (a *app) handelGatewayGetTaskResponse(task *corepb.Task_GatewayGetTaskResponse) {
	switch t := task.GatewayGetTaskResponse.GetResp().(type) {
	case *corepb.GatewayGetTaskResponse_Message:
		logg.L.Info("gateway get task from core success with 1 task")
		a._tasksAddToExecutor(t.Message, false)
	case *corepb.GatewayGetTaskResponse_Empty:
		logg.L.Info("gateway get task from core success but no task need to execute")
	}
}
