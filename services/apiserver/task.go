package apiserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/pkg/types/user"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
)

// taskID -> task([]bytes)

func CreateTask(taskID string, targetType task.Target, targetDeviceID string, taskBin []byte) error {
	dataCli := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)

	taskMq := dataCli.Mq
	e1 := taskMq.Publish(fmt.Sprintf("/task/%s/%s", targetType, targetDeviceID), []byte(taskID))     // 创建任务
	e2 := taskMq.Publish(fmt.Sprintf("/task/%s/%s/%s", targetType, targetDeviceID, taskID), taskBin) // 发送任务
	if errors.Join(e1, e2) != nil {
		return fmt.Errorf("create task error: %w", errors.Join(e1, e2))
	}

	var t task.Task
	switch targetType {
	case task.TargetGateway:
		t.Executor = user.RoleGateway
	case task.TargetAgent:
		t.Executor = user.RoleAgent
	case task.TargetCore:
		t.Executor = user.RoleCore
	}
	t.ExecutorID = targetDeviceID
	t.OperationCommend = string(taskBin)

	txn := dataCli.Begin()
	err := dataCli.CreateTask(context.Background(), txn, taskID, t)
	if err != nil {
		return fmt.Errorf("create task log error: %w", err)
	}
	dataCli.Commit(txn)

	return nil
}
