package services

import (
	"context"
	"fmt"
	"time"

	errCh "github.com/aenjoy/iot-lubricant/pkg/error"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/pkg/types/user"
	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"

	"github.com/bytedance/sonic"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

// _taskHelper 封装了任务相关的操作
// 返回：两个值，第一个是响应topic，第二个是taskID
func _taskHelper(
	ctx context.Context,
	txnHelper func() (*gorm.DB, *errCh.ErrorChan, func()),
	storeMq mq.Mq,
	dbAddAsyncJob func(context.Context, *gorm.DB, *model.AsyncJob) error,
	taskID *string,
	executorType user.Role,
	executorID string,
	userID string,
	taskName string,
	topicPrefix string,
	bin []byte,
) (string, string, error) {
	txn, errorCh, commit := txnHelper()
	defer commit()

	taskId := func() string {
		if taskID == nil {
			id := xid.New().String()
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
		UserID:    userID,
		Name:      taskName,
		Status:    "pending",
		Data:      taskString,
	})
	if err != nil {
		return "", "", err
	}

	topic := fmt.Sprintf("%s/%s", topicPrefix, executorID)
	logger.Debugf("send task %s to %s", taskId, topic)

	go func() {
		// 任务数据发送需要异步操作(在其它线程订阅这个topic后)，否则可能会导致获取任务失败
		pbTopic := fmt.Sprintf("%s/%s", topic, taskId)
		logger.Debugf("send task data to %s", pbTopic)
		time.Sleep(500 * time.Millisecond)
		_ = storeMq.PublishBytes(pbTopic, bin)
	}()

	err = storeMq.PublishBytes(topic, []byte(taskId))
	if err != nil {
		errorCh.Report(err, exceptionCode.MqPublishFailed, "failed to publish new task message to internal message queue", true)
		return "", "", err
	}
	return fmt.Sprintf("%s/%s/response", topic, taskId), taskId, nil
}
