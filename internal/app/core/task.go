package core

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/errs"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/user"
)

// taskID -> task([]bytes)
var (
	hasTask = sync.Map{} // targetID -> struct{} cache for get task
)

func CreateTask(taskID string, targetType task.Target, targetDeviceID string, taskBin []byte) error {
	dataCli := dataCli()
	hasTask.Store(targetDeviceID, struct{}{})
	taskMq := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore).Mq
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
func getTaskIDCh(ctx context.Context, targetType task.Target, targetDeviceID string) (chan string, error) {
	ch := make(chan string)
	taskMq := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore).Mq
	subscribe, err := taskMq.Subscribe(fmt.Sprintf("/task/%s/%s", targetType, targetDeviceID))
	if err != nil {
		return nil, err
	}

	go func() {
		select {
		case <-ctx.Done():
			break
		case taskID := <-subscribe:
			ch <- string(taskID)
		}
		close(ch)
	}()
	return ch, nil
}
func getTask(ctx context.Context, targetType task.Target, targetDeviceID, taskID string) ([]byte, error) {
	if _, ok := hasTask.Load(targetDeviceID); !ok {
		return nil, errs.ErrTargetNoTask
	}
	taskMq := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore).Mq
	t, err := taskMq.Subscribe(fmt.Sprintf("/task/%s/%s/%s", targetType, targetDeviceID, taskID))
	if err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, errs.ErrTimeout
	case task := <-t:
		hasTask.Delete(targetDeviceID)
		return task, nil
	}
}
