package model

import (
	"context"
	"errors"

	"github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"gorm.io/gorm"
)

var _ GatewayDbOperator = (*GatewayDb)(nil)
var _ CoreDbOperator = (*CoreDb)(nil)

type CoreDbOperator interface {
	Begin() *gorm.DB
	Commit(txn *gorm.DB)
	Rollback(txn *gorm.DB)

	// Common:
	GatewayIDGetUserID(ctx context.Context, id string) (string, error)
	AgentIDGetGatewayID(ctx context.Context, id string) (string, error)

	// Gateway:
	IsGatewayIdExists(string) bool
	GetGatewayInfo(ctx context.Context, id string) (*Gateway, error)
	AddGateway(ctx context.Context, txn *gorm.DB, userID string, gateway Gateway) error // need txn
	UpdateGateway(ctx context.Context, txn *gorm.DB, gateway Gateway) error             // need txn
	DeleteGateway(ctx context.Context, txn *gorm.DB, id string) error
	GetAllGatewayInfo(ctx context.Context) ([]Gateway, error)

	// Agent:
	GetAgentInfo(id string) (*Agent, error)
	AddAgent(ctx context.Context, txn *gorm.DB, gatewayID string, agent Agent) error
	UpdateAgent(ctx context.Context, txn *gorm.DB, agent Agent) error
	DeleteAgent(ctx context.Context, txn *gorm.DB, id string) error
	GetAgentList(ctx context.Context, gatewayID string) ([]Agent, error)

	// Data:
	StoreAgentGatherData(ctx context.Context, txn *gorm.DB, id, content string) error
	GetDataCleaner(id string) (*Clean, error)
	DeleteAgentGatherData(ctx context.Context, txn *gorm.DB, id string, timeStart int64, timeEnd int64) error

	// TaskLog:
	CreateTask(ctx context.Context, txn *gorm.DB, id string, task task.Task) error
	TaskUpdateOperationType(ctx context.Context, txn *gorm.DB, id string, Type task.Operation) error
	TaskUpdateOperationCommend(ctx context.Context, txn *gorm.DB, id string, commend string) error

	// User:
}

type GatewayDbOperator interface {
	GetServerInfo() ServerInfo
	IsAgentIdExists(string) bool
	GetAllAgentId() []string
	GetAllAgents() ([]Agent, error)
	RemoveAgent(...string) bool
	GetAgentReportCycle(string) int
	GetAgentGatherCycle(string) int
}

var (
	ErrNeedTxn = errors.New("this operation need start with txn support")
)
