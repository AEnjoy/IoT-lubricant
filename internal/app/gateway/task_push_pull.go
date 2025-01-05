package gateway

import (
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/google/uuid"
)

func (a *app) handelCorePushTaskAsync(task *core.Task_CorePushTaskRequest) *core.CorePushTaskResponse {
	var id string

	taskDetail := task.CorePushTaskRequest.GetMessage()
	if taskDetail.TaskId == "" {
		id = uuid.NewString()
		taskDetail.TaskId = id
	}

	a.task.AddTask(taskDetail)
	return &core.CorePushTaskResponse{
		TaskId: id,
	}
}
