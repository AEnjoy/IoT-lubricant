package edge

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type clientMqRecv struct {
	ctrl   context.Context
	cancel context.CancelFunc
	*chs
}
type chs struct {
	agentDevice      <-chan []byte
	agentDataPushAck <-chan []byte
	messagePushAck   <-chan []byte
	messagePull      <-chan []byte
}

func (c *clientMqRecv) SetContext(ctx context.Context) {
	c.ctrl, c.cancel = context.WithCancel(ctx)
}
func (c *clientMqRecv) handelCh() {
	for {
		select {
		case <-c.ctrl.Done():
			return
		case v := <-c.agentDevice: //客户端命令
			command := types.TaskCommand{}
			_ = json.Unmarshal(v, &command)
			if command.ID == task.OperationRemoveAgent {
				removeAgent()
			}
		}
	}
}

type clientMqSend struct {
	ctrl    context.Context
	cancel  context.CancelFunc
	agentID string
	mq      mq.Mq[[]byte]
}

func (c *clientMqSend) SetContext(ctx context.Context) {
	c.ctrl, c.cancel = context.WithCancel(ctx)
}

func (c *clientMqSend) send() error {
	for {
		select {
		case <-c.ctrl.Done():
			logger.Info("send data to gateway worker routine canceled")
			return nil
		case data := <-dataChan2:
			id := uuid.NewString()
			dataMessage := &gateway.DataMessage{
				Flag:      2,
				MessageId: id,
				AgentId:   c.agentID,
				Data:      data.Data,
				Time:      data.Timestamp.Format("2006-01-02 15:04:05"),
			}
			dataSend, err := json.Marshal(dataMessage)
			if err != nil {
				return err
			}

			err = c.mq.Publish(types.Topic_AgentDataPush+c.agentID, dataSend)
			if err != nil {
				return err
			}
		}
	}
}
