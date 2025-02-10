package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	exceptCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	taskTypes "github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/user"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

type AgentService struct {
	db    repo.CoreDbOperator
	store *datastore.DataStore
}

func (a *AgentService) PushTask(ctx context.Context, taskid *string, gatewayID, agentID string, bin []byte) (string, string, error) {
	txn, errorCh, commit := a.txnHelper()
	defer commit()

	mq := a.store.Mq
	taskId := func() string {
		if taskid == nil {
			id := uuid.NewString()
			taskid = &id
			return id
		}
		return *taskid
	}()
	task := taskTypes.Task{
		TaskID:   taskId,
		Operator: user.RoleCore,
	}
	taskString, _ := sonic.MarshalString(task)
	err := a.db.AddAsyncJob(ctx, txn, &model.AsyncJob{
		RequestID: taskId,
		Name:      "push task",
		Status:    "pending",
		Data:      taskString,
	})
	if err != nil {
		return "", "", err
	}

	topic := fmt.Sprintf("/task/%s/%s/%s/%s", taskTypes.TargetGateway, gatewayID, taskTypes.TargetAgent, agentID)
	err = errors.Join(
		mq.PublishBytes(fmt.Sprintf("%s", topic), []byte(taskId)),
		mq.PublishBytes(fmt.Sprintf("%s/%s", topic, taskId), bin))
	if err != nil {
		errorCh.Report(err, exceptCode.MqPublishFailed, "faild to publish new task message to internal message queue", true)
		return "", "", err
	}
	return fmt.Sprintf("%s/%s/response", topic, taskId), taskId, nil
}
