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
	GetServerInfo() *model.ServerInfo
	IsAgentIdExists(string) bool
	GetAllAgentId() []string
	GetAllAgents() ([]model.Agent, error)
	RemoveAgent(...string) bool
	GetAgentReportCycle(string) int
	GetAgentGatherCycle(string) int
	GetAgentInstance(id string) model.AgentInstance
}
