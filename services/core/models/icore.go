package models

import (
	"context"
	"time"

	model2 "github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type ICoreDb interface {
	Begin() *gorm.DB
	Commit(txn *gorm.DB)
	Rollback(txn *gorm.DB)

	// Common:
	GatewayIDGetUserID(ctx context.Context, id string) (string, error)
	AgentIDGetGatewayID(ctx context.Context, id string) (string, error)

	// Gateway:
	IsGatewayIdExists(id string) bool
	GetGatewayInfo(ctx context.Context, id string) (*model2.Gateway, error)
	AddGateway(ctx context.Context, txn *gorm.DB, userID string, gateway model2.Gateway) error // need txn
	UpdateGateway(ctx context.Context, txn *gorm.DB, gateway model2.Gateway) error             // need txn
	DeleteGateway(ctx context.Context, txn *gorm.DB, id string) error
	GetAllGatewayInfo(ctx context.Context) ([]model2.Gateway, error)
	AddGatewayHostInfo(ctx context.Context, txn *gorm.DB, info *model2.GatewayHost) error
	GetGatewayHostInfo(ctx context.Context, hostid string) (model2.GatewayHost, error)
	UpdateGatewayHostInfo(ctx context.Context, txn *gorm.DB, hostid string, info *model2.GatewayHost) error
	ListGatewayHostInfoByUserID(ctx context.Context, userID string) ([]model2.GatewayHost, error)
	DeleteGatewayHostInfo(ctx context.Context, txn *gorm.DB, hostId string) error

	// Agent:
	GetAgentInfo(id string) (*model2.Agent, error)
	AddAgent(ctx context.Context, txn *gorm.DB, gatewayID string, agent model2.Agent) error
	UpdateAgent(ctx context.Context, txn *gorm.DB, agent model2.Agent) error
	DeleteAgent(ctx context.Context, txn *gorm.DB, id string) error
	GetAgentList(ctx context.Context, gatewayID string) ([]model2.Agent, error)

	// Data:
	StoreAgentGatherData(ctx context.Context, txn *gorm.DB, id, content string) error
	GetDataCleaner(id string) (*model2.Clean, error)
	DeleteAgentGatherData(ctx context.Context, txn *gorm.DB, id string, timeStart int64, timeEnd int64) error

	// TaskLog:
	CreateTask(ctx context.Context, txn *gorm.DB, id string, task task.Task) error
	TaskUpdateOperationType(ctx context.Context, txn *gorm.DB, id string, Type task.Operation) error
	TaskUpdateOperationCommend(ctx context.Context, txn *gorm.DB, id string, commend string) error

	// ErrorLog:
	GetErrorLogs(ctx context.Context, gatewayid string, from, to time.Time, limit int) ([]model2.ErrorLogs, error)
	GetErrorLogByErrorID(ctx context.Context, errID string) (model2.ErrorLogs, error)

	// User:
	QueryUser(ctx context.Context, userName, uuid string) (model2.User, error)

	// Auth:
	SaveToken(ctx context.Context, tk *model2.Token) error
	SaveTokenOauth2(ctx context.Context, tk *oauth2.Token, userID string) error
	GetUserRefreshToken(ctx context.Context, userID string) (string, error)

	// Async Job
	AddAsyncJob(ctx context.Context, txn *gorm.DB, task *model2.AsyncJob) error
	GetAsyncJob(ctx context.Context, requestId string) (model2.AsyncJob, error)
	SetAsyncJobStatus(ctx context.Context, txn *gorm.DB, requestId string, status string) error

	// internal
	SetGatewayStatus(ctx context.Context, txn *gorm.DB, gatewayID, status string) error
	GetGatewayStatus(ctx context.Context, gatewayID string) (string, error)
}
