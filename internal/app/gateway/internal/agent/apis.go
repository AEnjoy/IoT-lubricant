package agent

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/pkg/docker"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
)

var _ Apis = (*agentApis)(nil)

type agentApis struct {
	db repo.GatewayDbOperator

	*pool
}

func (a *agentApis) init() {
	ctx := context.Background()
	agents, err := a.db.GetAllAgents()
	if err != nil {
		panic(err)
	}
	//todo: 这里可以加个并发限速
	for _, agent := range agents {
		go func(agent model.Agent) {
			ins := a.db.GetAgentInstance(agent.AgentId)

			control := new(agentControl)
			control.agentInfo = &agent

			if ins.Local && !docker.IsContainerRunning(ins.ContainerID) {
				if err := docker.StartContainer(ins.ContainerID); err != nil {
					logger.Error("start container failed", ins.ContainerID, err)
					return
				}
			}

			if err := a.JoinAgent(ctx, control); err != nil {
				logger.Error("agent join to handel pool failed", agent.AgentId, err)
			}
		}(agent)
	}
}
func (a *agentApis) StartAgent(id ...string) error {
	panic("implement me")
}
func (a *agentApis) StopAgent(id ...string) error {
	panic("implement me")
}
func (a *agentApis) KillAgent(id ...string) error {
	panic("implement me")
}
func (a *agentApis) RemoveAgent(id ...string) error {
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
	panic("implement me")
}
