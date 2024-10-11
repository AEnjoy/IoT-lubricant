package gateway

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
)

type agentCtrl struct {
	agentDevice <-chan []byte // /agentData/+agentID
	reg         <-chan []byte // Topic_AgentRegister
	ctx         context.Context
	ctrl        context.CancelFunc
}

type clientMq struct {
	ctrl       context.Context
	cancel     context.CancelFunc
	deviceList *sync.Map // agentId - agentCtrl channel
}

func (a *clientMq) Start() {

}
func (a *clientMq) Stop() {
	a.cancel()
}
func (a *clientMq) SetContext(ctx context.Context) {
	a.ctrl, a.cancel = context.WithCancel(ctx)
}

func (a *clientMq) handelGatewayInfo(in <-chan []byte, err error) error {
	if err != nil {
		return err
	}
	for {
		select {
		case <-a.ctrl.Done():
			return nil
		case <-in:
			// todo:handle the agent's request
		}
	}
}
func (a *clientMq) handelGatewayData(in <-chan []byte, err error) error {
	if err != nil {
		return err
	}
	for {
		select {
		case <-a.ctrl.Done():
			return nil
		case <-in:
			// todo:handle the agent's request
		}
	}
}
func (a *clientMq) handelPing(in <-chan []byte, err error) error {
	if err != nil {
		return err
	}
	for {
		select {
		case <-a.ctrl.Done():
			return nil
		case <-in:
			// todo:handle the agent's request
		}
	}
}
func (a *app) handelAgentDataPush(in <-chan []byte, err error, id string) error {
	if err != nil {
		return err
	}
	for {
		select {
		case <-a.ctrl.Done():
			return nil
		case data := <-in:
			var out gateway.DataMessage
			err := json.Unmarshal(data, &out)
			if err != nil {
				return err
			}
			if out.Flag == 2 {
				v, ok := agentStore.Load(id)
				if !ok {
					return ErrAgentNotFound
				}
				agentMap := v.(*agentData)
				agentMap.parseData(&out, a.GetAgentGatherCycle(id))
			}
			// todo: handle other flag
		}
	}
}
func (a *clientMq) handelAgentMessagePush(in <-chan []byte, err error, id string) error {
	if err != nil {
		return err
	}
	for {
		select {
		case <-a.ctrl.Done():
			return nil
		case <-in:
			// todo:handle the agent's request
		}
	}
}
