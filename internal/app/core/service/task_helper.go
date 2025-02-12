package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	errCh "github.com/AEnjoy/IoT-lubricant/pkg/error"
	exceptCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	taskTypes "github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/user"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func _taskHelper(
	ctx context.Context,
	txnHelper func() (*gorm.DB, *errCh.ErrorChan, func()),
	storeMq mq.Mq[[]byte],
	dbAddAsyncJob func(context.Context, *gorm.DB, *model.AsyncJob) error,
	taskID *string,
	executorType user.Role,
	executorID string,
	taskName string,
	topicPrefix string,
	bin []byte,
) (string, string, error) {
	txn, errorCh, commit := txnHelper()
	defer commit()

	taskId := func() string {
		if taskID == nil {
			id := uuid.NewString()
			taskID = &id
			return id
		}
		return *taskID
	}()

	task := taskTypes.Task{
		TaskID:     taskId,
		Operator:   user.RoleCore,
		Executor:   executorType,
		ExecutorID: executorID,
	}
	taskString, _ := sonic.MarshalString(task)
	err := dbAddAsyncJob(ctx, txn, &model.AsyncJob{
		RequestID: taskId,
		Name:      taskName,
		Status:    "pending",
		Data:      taskString,
	})
	if err != nil {
		return "", "", err
	}

	topic := fmt.Sprintf("%s/%s", topicPrefix, executorID)
	err = errors.Join(
		storeMq.PublishBytes(fmt.Sprintf("%s", topic), []byte(taskId)),
		storeMq.PublishBytes(fmt.Sprintf("%s/%s", topic, taskId), bin))
	if err != nil {
		errorCh.Report(err, exceptCode.MqPublishFailed, "failed to publish new task message to internal message queue", true)
		return "", "", err
	}
	return fmt.Sprintf("%s/%s/response", topic, taskId), taskId, nil
}
