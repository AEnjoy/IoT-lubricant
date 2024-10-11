package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
)

var (
	ErrAgentNotFound = errors.New("agentData not found")
)

func (a *app) handelAgentRegister(in <-chan []byte, err error) error {
	if err != nil {
		return err
	}
	for {
		select {
		case <-a.ctrl.Done():
			return nil
		case v := <-in:
			reg := &model.Register{}
			if err = json.Unmarshal(v, reg); err != nil {
				return err
			}
			ping, err := json.Marshal(model.Ping{Status: 1})
			if err != nil {
				return err
			}
			return a.mq.Publish(model.Topic_AgentRegisterAck+reg.ID, ping)
		}
	}
}

func (a *app) joinAgent(id string) (errs error) {
	ctx, cf := context.WithCancel(context.Background())
	ch := &agentCtrl{
		ctx:  ctx,
		ctrl: cf,
	}
	ag := &agentData{
		data:       make([]*gateway.DataMessage, 0),
		sendSignal: make(chan struct{}),
		l:          sync.Mutex{},
	}
	a.deviceList.Store(id, ch)
	agentStore.Store(id, ag)
	go func() {
		_ = a.pushDataToServer(ctx, id)
	}()

	go func() {
		chData, e := a.mq.Subscribe(model.Topic_AgentRegister + id)
		ch.reg = chData
		err := a.handelAgentRegister(chData, e)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}()
	go func() {
		err := a.subscribeDeviceMQ(ch, id)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}()
	go func() {
		ch, err := a.mq.Subscribe(model.Topic_AgentDataPush + id)
		err = a.handelAgentDataPush(ch, err, id)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}()
	go func() {
		ch, err := a.mq.Subscribe(model.Topic_MessagePush + id)
		err = a.handelAgentMessagePush(ch, err, id)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}()
	return
}

func (a *app) stopAgent(id string) (errs error) {
	v, ok := a.deviceList.Load(id)
	if !ok {
		return ErrAgentNotFound
	}
	ch := v.(*agentCtrl)
	ch.ctrl() // stop

	e1 := a.mq.Unsubscribe(model.Topic_AgentRegister+id, ch.reg)
	e2 := a.mq.Unsubscribe(model.Topic_AgentDevice+id, ch.agentDevice)
	e3 := a.mq.Unsubscribe(model.Topic_MessagePushAck+id, ch.messagePushAck)
	e4 := a.mq.Unsubscribe(model.Topic_MessagePull+id, ch.messagePull)

	commend, _ := json.Marshal(model.Command{ID: model.Command_RemoveAgent})
	data, _ := json.Marshal(gateway.DataMessage{Flag: 5, AgentId: id, Data: commend})
	e5 := a.mq.Publish(model.Topic_AgentDevice+id, data)

	//e5 := a.mq.Unsubscribe(model.Topic_AgentRegisterAck+id, ch.regAck)
	errs = errors.Join(errs, e1, e2, e3, e4, e5)

	a.deviceList.Delete(id)
	agentStore.Delete(id)
	a.GatewayDbCli.RemoveAgent(id)
	return
}
func (a *app) removeAgent(id ...string) bool {
	for _, s := range id {
		err := a.stopAgent(s)
		if err != nil {
			return false
		}
		// todo: remove agent data and other operation
	}
	return true
}

func (a *app) subscribeDeviceMQ(in *agentCtrl, id string) error {
	mq := a.mq
	in.agentDevice, _ = mq.Subscribe(model.Topic_AgentDevice + id)
	//in.regAck, _ = mq.Subscribe(model.Topic_AgentRegisterAck + id)
	in.messagePushAck, _ = mq.Subscribe(model.Topic_MessagePushAck + id)
	//in.messagePull, _ = mq.Subscribe(model.Topic_MessagePull + id)
	for {
		select {
		case <-a.ctrl.Done():
			return nil
		}
	}
}

func (a *app) initClientMq() (errs error) {
	mq := a.mq
	a.clientMq.deviceList = a.deviceList
	for _, id := range a.GatewayDbCli.GetAllAgentId() {
		if err := a.joinAgent(id); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	go func() {
		err := a.handelGatewayInfo(mq.Subscribe(model.Topic_GatewayInfo))
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}()
	go func() {
		err := a.handelGatewayData(mq.Subscribe(model.Topic_GatewayData))
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}()
	go func() {
		err := a.handelPing(mq.Subscribe(model.Topic_Ping))
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}()
	//go func() {
	//	err := a.handelAgentMessagePush(mq.Subscribe(model.Topic_MessagePush))
	//	if err != nil {
	//		errs = errors.Join(errs, err)
	//	}
	//}()

	return
}
