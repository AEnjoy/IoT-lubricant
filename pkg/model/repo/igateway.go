package repo

import (
	model2 "github.com/AEnjoy/IoT-lubricant/pkg/model"
	"gorm.io/gorm"
)

type IGatewayDb interface {
	Begin() *gorm.DB
	Commit(txn *gorm.DB)
	Rollback(txn *gorm.DB)

	AddOrUpdateServerInfo(txn *gorm.DB, info *model2.ServerInfo) error
	GetServerInfo(_ *gorm.DB) *model2.ServerInfo
	IsAgentIdExists(_ *gorm.DB, id string) bool
	GetAllAgentId(_ *gorm.DB) []string
	GetAllAgents(_ *gorm.DB) ([]model2.Agent, error)
	GetAgent(id string) (model2.Agent, error)
	UpdateAgent(txn *gorm.DB, id string, agent *model2.Agent) error
	RemoveAgent(txn *gorm.DB, id ...string) bool
	GetAgentReportCycle(_ *gorm.DB, id string) int
	GetAgentGatherCycle(_ *gorm.DB, id string) int
	GetAgentInstance(optionalTxn *gorm.DB, id string) model2.AgentInstance
	AddAgentInstance(txn *gorm.DB, ins *model2.AgentInstance) error
	UpdateAgentInstance(txn *gorm.DB, id string, ins model2.AgentInstance) error
}
