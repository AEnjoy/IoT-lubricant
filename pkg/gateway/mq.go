package gateway

import (
	"context"
	"sync"
)

type chs struct {
	agentDevice    <-chan []byte // /agent/+agentID
	regAck         <-chan []byte // /agent/register/ack/+agentID
	messagePushAck <-chan []byte // /gateway/message/push/ack/+agentID
	messagePull    <-chan []byte // /gateway/message/pull/+agentID
	reg            <-chan []byte // Topic_AgentRegister
}

type clientMq struct {
	ctrl       context.Context
	cancel     context.CancelFunc
	deviceList *sync.Map // agentId - chs channel
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
		}
	}
}
func (a *clientMq) handelAgentDataPush(in <-chan []byte, err error, id string) error {
	if err != nil {
		return err
	}
	for {
		select {
		case <-a.ctrl.Done():
			return nil
		}
	}
}
func (a *clientMq) handelAgentMessagePush(in <-chan []byte, err error) error {
	if err != nil {
		return err
	}
	for {
		select {
		case <-a.ctrl.Done():
			return nil
		}
	}
}
