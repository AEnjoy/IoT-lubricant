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

	var instance model.AgentInstance

	// 处理添加不在本机agent的情况
	if req.CreateAgentConf.AgentContainerInfo == nil && req.CreateAgentConf.DriverContainerInfo == nil &&
		req.AgentInfo.Address != "" {
		instance = model.AgentInstance{AgentId: req.AgentInfo.AgentId, IP: req.AgentInfo.Address, Online: true}
	} else {
		var driverIP, agentIP string
		if req.CreateAgentConf.DriverContainerInfo != nil {
			resp, err := docker.Create(context.Background(), req.CreateAgentConf.DriverContainerInfo)
			if err != nil {
				errorCh.Report(err, code2.AddAgentFailed, "add edge driver failed", true)
				return err
			}
			driverIP, err = docker.GetContainerIPAddress(context.Background(), resp.ID)
			if err != nil {
				errorCh.Report(err, code2.AddAgentFailed, "failed to get driver container ip", true)
				return err
			}
		}
		if req.CreateAgentConf.AgentContainerInfo != nil {
			req.CreateAgentConf.AgentContainerInfo.Env["DRIVER_IP"] = driverIP
			resp, err := docker.Create(context.Background(), req.CreateAgentConf.AgentContainerInfo)
			if err != nil {
				errorCh.Report(err, code2.AddAgentFailed, "add edge agent failed", true)
				return err
			}
			agentIP, err = docker.GetContainerIPAddress(context.Background(), resp.ID)
			if err != nil {
				errorCh.Report(err, code2.AddAgentFailed, "failed to get agent container ip", true)
				return err
			}
			instance.ContainerID = resp.ID
			instance.Local = true
			instance.AgentId = req.AgentInfo.AgentId
			instance.IP = fmt.Sprintf("%s:%d", agentIP, req.CreateAgentConf.AgentContainerInfo.ServicePort)
		}
	}

	err := a.db.AddAgentInstance(txn, instance)
	if err != nil {
		errorCh.Report(err, code2.AddAgentFailed, "add agent instance failed due to:%v", true, err)
		return err
	}
	err = a.pool.JoinAgent(context.Background(), newAgentControl(&req.AgentInfo))
	if err != nil {
		errorCh.Report(err, code2.AddAgentFailed, "add agent instance failed", true)
		return err
	}
	return nil
}

func (a *agentApis) isLocalAgentDevice(id string) bool {
	return a.db.GetAgentInstance(nil, id).Local
}
