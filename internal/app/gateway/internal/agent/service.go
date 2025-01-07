package agent

import (
	"sync"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
)

type Apis interface {
	StartAgent(id string) error
	StopAgent(id string) error
	KillAgent(id string) error
	RemoveAgent(id string) error
	UpdateAgent(id string) error
	EditAgent(id string, info model.Agent) error
	SetAgent(id string, info *agent.AgentInfo) error
	GetAgentInfo(id string) (*agent.AgentInfo, error)
	GetAgentModel(id string) (*model.Agent, error)
	AddAgent(req *model.CreateAgentRequest) error
	CreateAgent(req *model.CreateAgentRequest) error
}

var (
	_agentCli agent.EdgeServiceClient
	_apis     Apis
	once      sync.Once
)

// NewAgentApis 初始化并(或)获取AgentApis 对象  在后续获取对象时 可以设置参数为nil
func NewAgentApis(db repo.GatewayDbOperator) Apis {
	if (db == nil) && _apis == nil {
		panic("db and cli can not be nil at the init time")
	}

	if _apis == nil {
		once.Do(func() {
			a := &agentApis{db: db, pool: newPool()}
			go a.init()
			_apis = a
		})
	}
	return _apis
}
