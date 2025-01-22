package repo

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"gorm.io/gorm"
)

type CoreDbOperator interface {
	Begin() *gorm.DB
	Commit(txn *gorm.DB)
	Rollback(txn *gorm.DB)

	// Common:
	GatewayIDGetUserID(ctx context.Context, id string) (string, error)
	AgentIDGetGatewayID(ctx context.Context, id string) (string, error)

	// Gateway:
	IsGatewayIdExists(id string) bool
	GetGatewayInfo(ctx context.Context, id string) (*model.Gateway, error)
	AddGateway(ctx context.Context, txn *gorm.DB, userID string, gateway model.Gateway) error // need txn
	UpdateGateway(ctx context.Context, txn *gorm.DB, gateway model.Gateway) error             // need txn
	DeleteGateway(ctx context.Context, txn *gorm.DB, id string) error
	GetAllGatewayInfo(ctx context.Context) ([]model.Gateway, error)

	// Agent:
	GetAgentInfo(id string) (*model.Agent, error)
	AddAgent(ctx context.Context, txn *gorm.DB, gatewayID string, agent model.Agent) error
	UpdateAgent(ctx context.Context, txn *gorm.DB, agent model.Agent) error
	DeleteAgent(ctx context.Context, txn *gorm.DB, id string) error
	GetAgentList(ctx context.Context, gatewayID string) ([]model.Agent, error)

	// Data:
	StoreAgentGatherData(ctx context.Context, txn *gorm.DB, id, content string) error
	GetDataCleaner(id string) (*model.Clean, error)
	DeleteAgentGatherData(ctx context.Context, txn *gorm.DB, id string, timeStart int64, timeEnd int64) error

	// TaskLog:
	CreateTask(ctx context.Context, txn *gorm.DB, id string, task task.Task) error
	TaskUpdateOperationType(ctx context.Context, txn *gorm.DB, id string, Type task.Operation) error
	TaskUpdateOperationCommend(ctx context.Context, txn *gorm.DB, id string, commend string) error

	// User:
}
type GatewayDbOperator interface {
	Begin() *gorm.DB
	Commit(txn *gorm.DB)
	Rollback(txn *gorm.DB)

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
	AddAgentInstance(txn *gorm.DB, ins model.AgentInstance) error
	UpdateAgentInstance(txn *gorm.DB, id string, ins model.AgentInstance) error
}
