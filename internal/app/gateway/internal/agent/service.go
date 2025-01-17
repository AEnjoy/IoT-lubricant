package agent

import (
	"sync"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	agentpb "github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	proxypb "github.com/AEnjoy/IoT-lubricant/protobuf/proxy"
)

type Apis interface {
	StartAgent(id string) error
	StopAgent(id string) error
	KillAgent(id string) error
	RemoveAgent(id string) error
	UpdateAgent(id string, optionalConf *model.CreateAgentRequest) error
	EditAgent(id string, info *proxypb.EditAgentRequest) error
	SetAgent(id string, info *agentpb.AgentInfo) error
	GetAgentInfo(id string) (*agentpb.AgentInfo, error)
	GetAgentModel(id string) (*model.Agent, error)
	AddAgent(req *model.CreateAgentRequest) error
	CreateAgent(req *model.CreateAgentRequest) error
}

var (
	_apis Apis
	once  sync.Once
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
