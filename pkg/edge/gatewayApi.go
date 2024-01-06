package edge

import (
	"errors"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/goccy/go-json"
)

func (a *app) joinGateway() error {
	reg, err := json.Marshal(model.Register{ID: a.config.ID})
	if err != nil {
		return err
	}
	err = a.mq.Publish(model.Topic_AgentRegister+a.config.ID, reg)
	if err != nil {
		return err
	}
	ch, err := a.mq.Subscribe(model.Topic_AgentRegisterAck + a.config.ID)
	if err != nil {
		return err
	}
	var pong model.Ping
	err = json.Unmarshal(a.mq.GetPayLoad(ch), &pong)
	if err != nil {
		return err
	}
	if pong.Status != 1 {
		return errors.New(pong.Message)
	}
	return nil
}
func (a *app) initClientMq() (errs error) {
	mq := a.mq
	var ch = &chs{}
	ch1, e1 := mq.Subscribe(model.Topic_AgentDevice + a.config.ID)
	ch2, e2 := mq.Subscribe(model.Topic_MessagePushAck + a.config.ID)
	ch3, e3 := mq.Subscribe(model.Topic_MessagePull + a.config.ID)
	ch4, e4 := mq.Subscribe(model.Topic_AgentDataPushAck + a.config.ID)
	ch.agentDevice = ch1
	ch.messagePushAck = ch2
	ch.messagePull = ch3
	ch.agentDataPushAck = ch4
	a.clientMq.chs = ch
	return errors.Join(errs, e1, e2, e3, e4)
}
