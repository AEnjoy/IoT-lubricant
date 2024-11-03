package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
	"google.golang.org/grpc"
)

var (
	ErrAgentNotFound = errors.New("agent not found")
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
			reg := &types.Register{}
			if err = json.Unmarshal(v, reg); err != nil {
				return err
			}
			ping, err := json.Marshal(types.Ping{Status: 1})
			if err != nil {
				return err
			}
			return a.mq.Publish(types.Topic_AgentRegisterAck+reg.ID, ping)
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
		_ = a.handelSignal(id)
	}()
	go func() {
		_ = a.handelPushDataToServer(ctx, id)
	}()

	go func() {
		chData, e := a.mq.Subscribe(types.Topic_AgentRegister + id)
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
		ch, err := a.mq.Subscribe(types.Topic_AgentDataPush + id)
		if err != nil {
			errs = errors.Join(errs, err)
			return
		}
		err = a.handelAgentDataPush(ch, id)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}()
	go func() {
		ch, err := a.mq.Subscribe(types.Topic_MessagePush + id)
		if err != nil {
			errs = errors.Join(errs, err)
			return
		}
		err = a.handelAgentMessagePush(ch, id)
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

	e1 := a.mq.Unsubscribe(types.Topic_AgentRegister+id, ch.reg)
	e2 := a.mq.Unsubscribe(types.Topic_AgentDevice+id, ch.agentDevice)

	commend, _ := json.Marshal(types.TaskCommand{ID: task.OperationRemoveAgent})
	data, _ := json.Marshal(gateway.DataMessage{Flag: 5, AgentId: id, Data: commend})
	e3 := a.mq.Publish(types.Topic_AgentDevice+id, data)

	//e5 := a.mq.Unsubscribe(model.Topic_AgentRegisterAck+id, ch.regAck)
	errs = errors.Join(errs, e1, e2, e3)

	a.deviceList.Delete(id)
	agentStore.Delete(id)
	a.GatewayDbOperator.RemoveAgent(id)
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
	in.agentDevice, _ = mq.Subscribe(types.Topic_AgentDevice + id)

	for {
		select {
		case <-a.ctrl.Done():
			return nil
		case <-in.agentDevice:
			// todo: handle agent command
		}
	}
}

func (a *app) initClientMq() (errs error) {
	mq := a.mq
	for _, id := range a.GatewayDbOperator.GetAllAgentId() {
		if err := a.joinAgent(id); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	go func() {
		err := a.handelGatewayInfo(mq.Subscribe(types.Topic_GatewayInfo))
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}()
	go func() {
		err := a.handelGatewayData(mq.Subscribe(types.Topic_GatewayData))
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}()
	go func() {
		err := a.handelPing(mq.Subscribe(types.Topic_Ping))
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}()
	// todo: handle message push
	//go func() {
	//	err := a.handelAgentMessagePush(mq.Subscribe(model.Topic_MessagePush))
	//	if err != nil {
	//		errs = errors.Join(errs, err)
	//	}
	//}()

	return
}
func (a *app) initAgentPool() error {
	as, err := a.GetAllAgents()
	if err != nil {
		return err
	}
	var errs error

	for _, agent := range as {
		control := new(agentControl)
		control.agentInfo = &agent
		go func() {
			err := control.joinAgentPool(a.ctrl)
			if err != nil {
				errs = errors.Join(errs, err)
			}
		}()
	}
	return errs
}

var agentPool = make(map[string]*agentControl)
var agentPoolMutex sync.Mutex

type agentControl struct {
	id        string
	agentInfo *types.Agent
	agentCli  agent.EdgeServiceClient
	ctx       context.Context
	cancel    context.CancelFunc
	gather    bool
}

func (a *agentControl) joinAgentPool(ctx context.Context) error {
	a.ctx, a.cancel = context.WithCancel(ctx)
	conn, err := grpc.NewClient(a.agentInfo.Address)
	if err != nil {
		return err
	}
	a.id = a.agentInfo.Id
	a.agentCli = agent.NewEdgeServiceClient(conn)

	agentPoolMutex.Lock()
	defer agentPoolMutex.Unlock()

	agentPool[a.id] = a
	return nil
}

func (a *agentControl) gatherData() {
	a.gather = true
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-time.Tick(time.Second * time.Duration(a.agentInfo.Cycle)):
			gatherDataResp, err := a.agentCli.Data(a.ctx, &agent.GetDataRequest{AgentID: a.id})
			if err != nil {
				logger.Errorf("Get Agent data failed: due to `%s`, agentID is `%s`", err.Error(), a.id)
				continue
			}
			gatherDataResp.GetData()
		}
	}
}
