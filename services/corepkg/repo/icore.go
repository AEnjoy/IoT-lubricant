package repo

import (
	"context"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/types/operation"
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
	IsGatewayIdExists(userID, gatewayID string) bool
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
	GetAgentsByProjectID(ctx context.Context, txn *gorm.DB, projectID string) ([]model.Agent, error)
	AddAgent(ctx context.Context, txn *gorm.DB, gatewayID string, agent model.Agent) error
	UpdateAgent(ctx context.Context, txn *gorm.DB, agent model.Agent) error
	UpdateAgentStatus(ctx context.Context, txn *gorm.DB, agentID, status string) error
	GetAgentStatus(ctx context.Context, agentID string) (string, error)
	DeleteAgent(ctx context.Context, txn *gorm.DB, id string) error
	GetAgentList(ctx context.Context, userID, gatewayID string) ([]model.Agent, error)
	GetAgentIDByAgentNameAndUserID(ctx context.Context, agentName, userID string) (string, error)

	// Data:
	StoreAgentGatherData(ctx context.Context, txn *gorm.DB, id, content string) error
	GetDataCleaner(id string) (*model.Clean, error)
	DeleteAgentGatherData(ctx context.Context, txn *gorm.DB, id string, timeStart int64, timeEnd int64) error

	// TaskLog:
	CreateTask(ctx context.Context, txn *gorm.DB, id string, task task.Task) error
	TaskUpdateOperationType(ctx context.Context, txn *gorm.DB, id string, Type operation.Operation) error
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
	GetAsyncJobResult(ctx context.Context, requestId string) (status, result string, err error)
	UserGetAsyncJobs(ctx context.Context, userID string, current, limit int) ([]model.AsyncJob, error)
	// SetAsyncJobStatus txn is allowed set to nil(means no txn)
	SetAsyncJobStatus(ctx context.Context, txn *gorm.DB, requestId string, status, result string) error

	// internal
	SetGatewayStatus(ctx context.Context, txn *gorm.DB, userid, gatewayID, status string) error
	GetGatewayStatus(ctx context.Context, gatewayID string) (string, error)

	// project
	AddProject(ctx context.Context, txn *gorm.DB, userid, projectid, projectname, description string) error
	RemoveProject(ctx context.Context, txn *gorm.DB, projectid string) error
	GetProject(ctx context.Context, projectid string) (model.Project, error)
	ListProject(ctx context.Context, userID string) ([]model.Project, error)

	GetProjectAgentNumber(ctx context.Context, txn *gorm.DB, projectid string) (int, error)
	AddDataStoreEngine(ctx context.Context, txn *gorm.DB, projectid, dsn, dataBaseType, description string) error
	UpdateEngineInfo(ctx context.Context, txn *gorm.DB, projectid, dsn, dataBaseType, description string) error
	GetEngineByProjectID(ctx context.Context, projectid string) (model.DataStoreEngine, error)
	BindProject(ctx context.Context, txn *gorm.DB, projectid string, agents []string) error
	AddWasher(ctx context.Context, txn *gorm.DB, w *model.Clean) error
	BindWasher(ctx context.Context, txn *gorm.DB, projectid string, washerID int, agentIDs []string) error
}
