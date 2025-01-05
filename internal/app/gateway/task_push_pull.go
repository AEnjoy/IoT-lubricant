package gateway

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/google/uuid"
)

func (a *app) _tasksAddToExecutor(taskDetail *core.TaskDetail, notice bool) *core.CorePushTaskResponse {
	var id string
	if taskDetail.TaskId == "" {
		id = uuid.NewString()
		taskDetail.TaskId = id
	}

	a.task.AddTask(taskDetail, notice)
	return &core.CorePushTaskResponse{
		TaskId: id,
	}
}
func (a *app) handelCorePushTaskAsync(task *core.Task_CorePushTaskRequest) *core.CorePushTaskResponse {
	return a._tasksAddToExecutor(task.CorePushTaskRequest.GetMessage(), true)
}

func (a *app) handelGatewayGetTaskResponse(task *core.Task_GatewayGetTaskResponse) {
	switch t := task.GatewayGetTaskResponse.GetResp().(type) {
	case *core.GatewayGetTaskResponse_Message:
		logger.Infoln("gateway get task from core success with 1 task")
		a._tasksAddToExecutor(t.Message, false)
	case *core.GatewayGetTaskResponse_Empty:
		logger.Infoln("gateway get task from core success but no task need to execute")
	}
}
