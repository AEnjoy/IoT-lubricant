package repo

import (
	"context"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/types/task"

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
	GetGatewayInfo(ctx context.Context, id string) (*model.Gateway, error)
	AddGateway(ctx context.Context, txn *gorm.DB, userID string, gateway model.Gateway) error // need txn
	UpdateGateway(ctx context.Context, txn *gorm.DB, gateway model.Gateway) error             // need txn
	DeleteGateway(ctx context.Context, txn *gorm.DB, id string) error
	GetAllGatewayInfo(ctx context.Context) ([]model.Gateway, error)
	GetAllGatewayByUserID(ctx context.Context, userID string) ([]model.Gateway, error)
	AddGatewayHostInfo(ctx context.Context, txn *gorm.DB, info *model.GatewayHost) error
	GetGatewayHostInfo(ctx context.Context, hostid string) (model.GatewayHost, error)
	UpdateGatewayHostInfo(ctx context.Context, txn *gorm.DB, hostid string, info *model.GatewayHost) error
	ListGatewayHostInfoByUserID(ctx context.Context, userID string) ([]model.GatewayHost, error)
	DeleteGatewayHostInfo(ctx context.Context, txn *gorm.DB, hostId string) error

	// Agent:
	GetAgentInfo(id string) (*model.Agent, error)
	AddAgent(ctx context.Context, txn *gorm.DB, gatewayID string, agent model.Agent) error
	UpdateAgent(ctx context.Context, txn *gorm.DB, agent model.Agent) error
	UpdateAgentStatus(ctx context.Context, txn *gorm.DB, agentID, status string) error
	GetAgentStatus(ctx context.Context, agentID string) (string, error)
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

	// ErrorLog:
	GetErrorLogs(ctx context.Context, gatewayid string, from, to time.Time, limit int) ([]model.ErrorLogs, error)
	GetErrorLogByErrorID(ctx context.Context, errID string) (model.ErrorLogs, error)
	SaveErrorLog(ctx context.Context, err *model.ErrorLogs) error

	// User:
	QueryUser(ctx context.Context, userName, uuid string) (model.User, error)

	// Auth:
	SaveToken(ctx context.Context, tk *model.Token) error
	SaveTokenOauth2(ctx context.Context, tk *oauth2.Token, userID string) error
	GetUserRefreshToken(ctx context.Context, userID string) (string, error)

	// Async Job
	AddAsyncJob(ctx context.Context, txn *gorm.DB, task *model.AsyncJob) error
	GetAsyncJob(ctx context.Context, requestId string) (model.AsyncJob, error)
	//GetAsyncJob(ctx context.Context, requestId string) (model.AsyncJob, error)
	SetAsyncJobStatus(ctx context.Context, txn *gorm.DB, requestId string, status, result string) error

	// internal
	SetGatewayStatus(ctx context.Context, txn *gorm.DB, gatewayID, status string) error
	GetGatewayStatus(ctx context.Context, gatewayID string) (string, error)
}
