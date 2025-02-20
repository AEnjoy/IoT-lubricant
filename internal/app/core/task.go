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

func CreateTask(taskID string, targetType task.Target, targetDeviceID string, taskBin []byte) error {
	dataCli := dataCli()

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

var taskChMap sync.Map

func getTaskIDCh(ctx context.Context, targetType task.Target, targetDeviceID string) (chan string, func(), error) {
	topic := fmt.Sprintf("/task/%s/%s", targetType, targetDeviceID)

	if val, exists := taskChMap.Load(topic); exists {
		entry := val.(struct {
			ch     chan string
			cancel func()
		})
		return entry.ch, entry.cancel, nil
	}

	var closeOnce sync.Once
	ch := make(chan string)
	taskMq := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore).Mq

	subscribe, err := taskMq.Subscribe(topic)
	if err != nil {
		return nil, nil, err
	}

	_clean := func() {
		_ = taskMq.Unsubscribe(topic, subscribe)
		taskChMap.Delete(topic)
		close(ch)
	}
	cancel := func() {
		closeOnce.Do(_clean)
	}

	// 存储到map（使用匿名结构体存储组合值）
	taskChMap.Store(topic, struct {
		ch     chan string
		cancel func()
	}{ch, cancel})

	go func() {
		defer closeOnce.Do(_clean)
		for {
			select {
			case <-ctx.Done():
				return
			case taskID := <-subscribe:
				ch <- string(taskID)
				//if id := string(taskID); id != "" {
				//	ch <- id
				//}
			}
		}
	}()

	return ch, cancel, nil
}

func getTask(ctx context.Context, targetType task.Target, targetDeviceID, taskID string) ([]byte, error) {
	taskMq := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore).Mq
	topic := fmt.Sprintf("/task/%s/%s/%s", targetType, targetDeviceID, taskID)
	t, err := taskMq.Subscribe(topic)
	if err != nil {
		return nil, err
	}

	ctx, cancel := createTimeOutContext(ctx)
	defer func() {
		cancel()
		_ = taskMq.Unsubscribe(topic, t)
	}()

	select {
	case <-ctx.Done():
		return nil, errs.ErrTargetNoTask
	case task := <-t:
		return task, nil
	}
}
