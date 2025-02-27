package repo

import (
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"gorm.io/gorm"
)

type IGatewayDb interface {
	Begin() *gorm.DB
	Commit(txn *gorm.DB)
	Rollback(txn *gorm.DB)

	AddOrUpdateServerInfo(txn *gorm.DB, info *model.ServerInfo) error
	GetServerInfo(_ *gorm.DB) *model.ServerInfo
	IsAgentIdExists(_ *gorm.DB, id string) bool
	GetAllAgentId(_ *gorm.DB) []string
	GetAllAgents(_ *gorm.DB) ([]model.Agent, error)
	GetAgent(id string) (model.Agent, error)
	UpdateAgent(txn *gorm.DB, id string, agent *model.Agent) error
	RemoveAgent(txn *gorm.DB, id ...string) bool
	GetAgentReportCycle(_ *gorm.DB, id string) int
	GetAgentGatherCycle(_ *gorm.DB, id string) int
	GetAgentInstance(optionalTxn *gorm.DB, id string) model.AgentInstance
	AddAgentInstance(txn *gorm.DB, ins *model.AgentInstance) error
	UpdateAgentInstance(txn *gorm.DB, id string, ins model.AgentInstance) error
}
