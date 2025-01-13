package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/pkg/docker"
	errCh "github.com/AEnjoy/IoT-lubricant/pkg/error"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/code"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	code2 "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	except "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
)

var _ Apis = (*agentApis)(nil)

type agentApis struct {
	db repo.GatewayDbOperator

	*pool
}

func (a *agentApis) init() {
	ctx := context.Background()
	agents, err := a.db.GetAllAgents(nil)
	if err != nil {
		panic(err)
	}
	//todo: 这里可以加个并发限速
	for _, agent := range agents {
		go func(agent model.Agent) {
			ins := a.db.GetAgentInstance(nil, agent.AgentId)
			control := newAgentControl(&agent)

			if ins.Local && !docker.IsContainerRunning(ins.ContainerID) {
				if err := docker.StartContainer(ins.ContainerID); err != nil {
					logger.Error("start container failed", ins.ContainerID, err)
					return
				}
			}

			if err := a.pool.JoinAgent(ctx, control); err != nil {
				logger.Error("agent join to handel pool failed", agent.AgentId, err)
			}
		}(agent)
		time.Sleep(100 * time.Millisecond) //避免Goroutine启动过快失败
	}
}
func (a *agentApis) StartAgent(id string) error {
	ctrl := a.pool.GetAgentControl(id)
	if ctrl == nil {
		return exception.New(except.ErrorGatewayAgentNotFound, exception.WithLevel(code.Error), exception.WithMsg(fmt.Sprintf("agentID:%s", id)))
	}
	if !a.isLocalAgentDevice(id) {
		return exception.New(except.OperationOnlyAtLocal, exception.WithLevel(code.Error),
			exception.WithMsg(fmt.Sprintf("agentID:%s", id)),
			exception.WithReason("Cannot control manually added agent instances to start, please manage manually."))
	}

	ins := a.db.GetAgentInstance(nil, id)
	if !docker.IsContainerRunning(ins.ContainerID) {
		err := bootAgentInstance(ins.ContainerID)
		if err != nil {
			return exception.ErrNewException(err, except.ErrorAgentStartFailed,
				exception.WithLevel(code.Error),
				exception.WithMsg(fmt.Sprintf("agentID:%s", id)))
		}
		return nil
	}
	logger.Warnln("agent already started", id)
	return nil
}
func (a *agentApis) StopAgent(id string) error {
	panic("implement me")
}
func (a *agentApis) KillAgent(id string) error {
	panic("implement me")
}
func (a *agentApis) RemoveAgent(id string) error {
	panic("implement me")
}
func (a *agentApis) UpdateAgent(id string) error {
	panic("implement me")
}
func (a *agentApis) EditAgent(id string, info model.Agent) error {
	panic("implement me")
}
func (a *agentApis) SetAgent(id string, info *agent.AgentInfo) error {
	panic("implement me")
}
func (a *agentApis) GetAgentInfo(id string) (*agent.AgentInfo, error) {
	panic("implement me")
}
func (a *agentApis) GetAgentModel(id string) (*model.Agent, error) {
	panic("implement me")
}
func (a *agentApis) AddAgent(req *model.CreateAgentRequest) error {
	return a.CreateAgent(req)
}
func (a *agentApis) CreateAgent(req *model.CreateAgentRequest) error {
	txn := a.db.Begin()
	errorCh := errCh.NewErrorChan()
	defer errCh.HandleErrorCh(errorCh).ErrorWillDo(func(error) {
		a.db.Rollback(txn)
	}).SuccessWillDo(func() {
		a.db.Commit(txn)
	}).Do()
	// 处理添加不在本机agent的情况
	if req.CreateAgentConf.AgentContainerInfo == nil && req.CreateAgentConf.DriverContainerInfo == nil &&
		req.AgentInfo.Address != "" {
		err := a.db.AddAgentInstance(txn, model.AgentInstance{AgentId: req.AgentInfo.AgentId, IP: req.AgentInfo.Address, Online: true})
		if err != nil {
			errorCh.Report(err, code2.AddAgentFailed, "add agent instance failed", true)
			return err
		}
		err = a.pool.JoinAgent(context.Background(), newAgentControl(&req.AgentInfo))
		if err != nil {
			errorCh.Report(err, code2.AddAgentFailed, "add agent instance failed", true)
			return err
		}
	} else {

	}
	panic("implement me")
}

func (a *agentApis) isLocalAgentDevice(id string) bool {
	return a.db.GetAgentInstance(nil, id).Local
}
