package edge

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/goccy/go-json"
)

type clientMq struct {
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

func (c *clientMq) SetContext(ctx context.Context) {
	c.ctrl, c.cancel = context.WithCancel(ctx)
}
func (c *clientMq) handelCh() {
	for {
		select {
		case <-c.ctrl.Done():
			return
		case v := <-c.agentDevice:
			command := model.Command{}
			_ = json.Unmarshal(v, &command)
			if command.ID == model.Command_RemoveAgent {
				removeAgent()
			}
		}
	}
}
