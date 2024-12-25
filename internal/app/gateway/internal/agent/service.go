package agent

import (
	"sync"

	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
)

type Apis interface {
}

var (
	_agentCli agent.EdgeServiceClient
	_apis     Apis
	once      sync.Once
)

// NewAgentApis 初始化并(或)获取AgentApis 对象  在后续获取对象时 可以设置参数为nil
func NewAgentApis(db repo.GatewayDbOperator, cli agent.EdgeServiceClient) Apis {
	if (db == nil || cli == nil) && _apis == nil {
		panic("db and cli can not be nil at the init time")
	}

	if _apis == nil {
		once.Do(func() {
			a := &agentApis{db: db, pool: newPool()}
			_agentCli = cli
			go a.init()
			_apis = a
		})
	}
	return _apis
}
