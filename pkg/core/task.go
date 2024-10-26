package core

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/google/uuid"
)

var (
	ErrTargetNoTask = errors.New("target has no task")
	ErrTimeout      = errors.New("get task timeout")
)

// taskID -> task([]bytes)
var (
	taskMq  = mq.NewMq[[]byte]()
	hasTask = sync.Map{} // targetID -> struct{} cache for get task
)

func createTask(targetType, targetDeviceID string, task []byte) error {
	taskID := uuid.NewString()
	hasTask.Store(targetDeviceID, struct{}{})
	e1 := taskMq.Publish(fmt.Sprintf("/task/%s/%s", targetType, targetDeviceID), []byte(taskID))  // 创建任务
	e2 := taskMq.Publish(fmt.Sprintf("/task/%s/%s/%s", targetType, targetDeviceID, taskID), task) // 发送任务
	return errors.Join(e1, e2)
}
func getTaskIDCh(ctx context.Context, targetType, targetDeviceID string) (chan<- string, error) {
	ch := make(chan string)
	subscribe, err := taskMq.Subscribe(fmt.Sprintf("/task/%s/%s", targetType, targetDeviceID))
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case taskID := <-subscribe:
				ch <- string(taskID)
			}
		}
	}()
	return ch, nil
}
func getTask(ctx context.Context, targetType, targetDeviceID, taskID string) ([]byte, error) {
	if _, ok := hasTask.Load(targetDeviceID); !ok {
		return nil, ErrTargetNoTask
	}

	t, err := taskMq.Subscribe(fmt.Sprintf("/task/%s/%s/%s", targetType, targetDeviceID, taskID))
	if err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrTimeout
	case task := <-t:
		hasTask.Delete(targetDeviceID)
		return task, nil
	}
}
