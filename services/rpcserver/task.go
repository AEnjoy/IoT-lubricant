package rpcserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/types/errs"
	"github.com/aenjoy/iot-lubricant/pkg/types/task"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
)

var taskChMap sync.Map

func (i *PbCoreServiceImpl) getTaskIDCh(ctx context.Context, targetType task.Target, userid, targetDeviceID string) (chan string, func(), error) {
	topic := fmt.Sprintf("/task/%s/%s/%s", userid, targetType, targetDeviceID)
	logger.Debugf("get task id chan topic： %s", topic)

	if val, exists := taskChMap.Load(topic); exists {
		entry := val.(struct {
			ch     chan string
			cancel func()
		})
		return entry.ch, entry.cancel, nil
	}

	var closeOnce sync.Once
	ch := make(chan string)
	taskMq := i.DataStore.Mq

	subscribe, err := taskMq.SubscribeBytes(topic)
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
				if taskID == nil {
					logg.L.Error("failed to get taskid from mq", "taskID is nil")
				} else {
					logg.L.Debugf("%v", taskID)
					ch <- string(taskID.([]byte))
				}

				//if id := string(taskID); id != "" {
				//	ch <- id
				//}
			}
		}
	}()

	return ch, cancel, nil
}

func (i *PbCoreServiceImpl) getTask(_ context.Context, targetType task.Target, userid, targetDeviceID, taskID string) ([]byte, error) {
	taskMq := i.DataStore.Mq
	topic := fmt.Sprintf("/task/%s/%s/%s/%s", userid, targetType, targetDeviceID, taskID)
	logg.L.Debugf("get task topic： %s", topic)
	t, err := taskMq.SubscribeBytes(topic)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = taskMq.Unsubscribe(topic, t)
	}()

	select {
	case <-time.After(3 * time.Second):
		return nil, errs.ErrTargetNoTask
	case task := <-t:
		return task.([]byte), nil
	}
}
