package gateway

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
)

var (
	ErrAgentNotFound = errors.New("agent not found")
)

func (a *app) dataStoreJoinAgent(id string) (errs error) {
	ag := &agentData{
		data:       make([]*agent.DataMessage, 0),
		sendSignal: make(chan struct{}),
		l:          sync.Mutex{},
	}
	agentDataStore.Store(id, ag)
	return
}

func (a *app) agentStop(id string) (errs error) {
	if _, ok := agentPool[id]; !ok {
		return ErrAgentNotFound
	}
	agentPool[id].cancel()
	return
}
func (a *app) agentRemove(id ...string) bool {
	for _, s := range id {
		err := a.agentStop(s)
		if err != nil {
			return false
		}

		agentDataStore.Delete(id)
		a.GatewayDbOperator.RemoveAgent(s)
		// todo: remove agent data and other operation
	}
	return true
}
func (a *app) agentPoolInit() error {
	as, err := a.GetAllAgents()
	if err != nil {
		return err
	}
	var errs error
	var wg sync.WaitGroup
	for _, agent := range as {
		wg.Add(1)
		_ = a.dataStoreJoinAgent(agent.Id)
		control := new(agentControl)
		control.agentInfo = &agent
		go func() {
			defer wg.Done()
			err := control.agentJoinToPool(a.ctrl)
			if err != nil {
				errs = errors.Join(errs, err)
			}
		}()
	}
	wg.Wait()
	return errs
}

var (
	agentPool      = make(map[string]*agentControl)
	agentPoolMutex sync.Mutex
	agentPoolCh    = make(chan *agentControl)
	joinSignal     = make(chan struct{})
)

type agentControl struct {
	id string
	// slot []int // for api paths

	agentInfo *model.Agent
	agentCli  agent.EdgeServiceClient
	ctx       context.Context
	cancel    context.CancelFunc
	gather    bool
	start     bool
}

func (a *agentControl) agentJoinToPool(ctx context.Context) error {
	a.ctx, a.cancel = context.WithCancel(ctx)

	cli, err := edge.NewAgentCli(a.agentInfo.Address)
	if err != nil {
		return err
	}

	a.id = a.agentInfo.Id
	a.agentCli = cli

	agentPoolMutex.Lock()
	defer agentPoolMutex.Unlock()

	agentPool[a.id] = a
	return nil
}
func (a *app) agentPoolAgentRegis() {
	agentPoolMutex.Lock()
	for _, a := range agentPool {
		if !a.start {
			agentPoolCh <- a
			a.start = true
		}
	}
	agentPoolMutex.Unlock()

	for _ = range joinSignal {
		agentPoolMutex.Lock()
		for _, v := range agentPool {
			if !v.start {
				agentPoolCh <- v
			}
		}
		agentPoolMutex.Unlock()
	}
}
func (a *app) agentPoolChStartService() {
	for {
		select {
		case <-a.ctrl.Done():
			return
		case control := <-agentPoolCh:
			control.start = true
			go control.gatherData()
			go a.agentHandelSignal(control.id)
		}
	}
}

func (a *agentControl) gatherData() {
	a.gather = true
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-time.Tick(time.Second * time.Duration(a.agentInfo.Cycle)):
			gatherDataResp, err := a.agentCli.GetGatherData(a.ctx, &agent.GetDataRequest{AgentID: a.id})
			if err != nil {
				logger.Errorf("Get Agent data failed: due to `%s`, agentID is `%s`", err.Error(), a.id)
				continue
			}
			gatherDataResp.GetData()
			if v, ok := agentDataStore.Load(a.id); ok {
				v.(*agentData).parseData(gatherDataResp)
			}
		}
	}
}
