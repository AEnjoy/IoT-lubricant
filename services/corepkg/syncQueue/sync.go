package syncQueue

import (
	"fmt"
	"time"

	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"

	"google.golang.org/protobuf/proto"
)

// SyncTaskQueue
//
//	对于所有的同步任务请求，都使用这个对象进行处理
type SyncTaskQueue struct {
	mq.Mq
}

func (s *SyncTaskQueue) Init() error {
	s.Mq = ioc.Controller.Get(ioc.APP_NAME_CORE_Internal_MQ_SERVICE).(mq.Mq)
	return nil
}

func (SyncTaskQueue) Weight() uint16 {
	return ioc.CoreSyncTaskSystem
}

func (SyncTaskQueue) Version() string {
	return ""
}

func (s *SyncTaskQueue) WaitTask(taskid string, timeout time.Duration) (*corepb.QueryTaskResultResponse, error) {
	var retVal *corepb.QueryTaskResultResponse
	ch, err := s.Mq.SubscribeBytes(fmt.Sprintf("/task/%s/sync/%s", taskTypes.TargetCore, taskid))
	if err != nil {
		return nil, err
	}
	defer func(Mq mq.Mq, topic string, sub <-chan any) {
		err := Mq.Unsubscribe(topic, sub)
		if err != nil {
			logg.L.Errorf("failed to unsubscribe from message queue: %v", err)
		}
	}(s.Mq, fmt.Sprintf("/task/%s/sync/%s", taskTypes.TargetCore, taskid), ch)

	select {
	case data := <-ch:
		retVal = &corepb.QueryTaskResultResponse{}
		err = proto.Unmarshal(data.([]byte), retVal)
		if err != nil {
			return nil, err
		}
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout")
	}
	return retVal, nil
}
func (s *SyncTaskQueue) FinshTask(taskid string, result *corepb.QueryTaskResultResponse) error {
	data, err := proto.Marshal(result)
	if err != nil {
		return err
	}
	return s.Mq.PublishBytes(fmt.Sprintf("/task/%s/sync/%s", taskTypes.TargetCore, taskid), data)
}
