package repo

import (
	"context"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/rs/xid"
	"golang.org/x/oauth2"

	"github.com/aenjoy/iot-lubricant/pkg/types/errs"
	"github.com/aenjoy/iot-lubricant/pkg/types/task"
	"gorm.io/gorm"
)

var _ ICoreDb = (*CoreDb)(nil)

type CoreDb struct {
	db *gorm.DB
}

func (d *CoreDb) GetAgentStatus(ctx context.Context, agentID string) (string, error) {
	var ag model.Agent
	err := d.db.WithContext(ctx).Model(model.Agent{}).Where("agent_id = ?", agentID).First(&ag).Error
	return ag.Status, err
}

func (d *CoreDb) UpdateAgentStatus(ctx context.Context, txn *gorm.DB, agentID, status string) error {
	return txn.WithContext(ctx).Model(model.Agent{}).
		Where("agent_id = ?", agentID).
		Update("status", status).Update("updated_at", time.Now()).Error
}

func (d *CoreDb) SaveErrorLog(ctx context.Context, err *model.ErrorLogs) error {
	err.ErrID = xid.New().String()
	err.CreatedAt = time.Now()
	return d.db.WithContext(ctx).Model(model.ErrorLogs{}).Create(err).Error
}

func (d *CoreDb) GetAllGatewayByUserID(ctx context.Context, userID string) ([]model.Gateway, error) {
	var ret []model.Gateway
	err := d.db.WithContext(ctx).Model(model.Gateway{}).Where("user_id = ?", userID).Find(&ret).Error
	return ret, err
}

func (d *CoreDb) GetGatewayStatus(ctx context.Context, gatewayID string) (string, error) {
	var ret model.Gateway
	return ret.Status, d.db.WithContext(ctx).Where("gateway_id = ?", gatewayID).First(&ret).Error
}

func (d *CoreDb) SetGatewayStatus(ctx context.Context, txn *gorm.DB, gatewayID, status string) error {
	return txn.WithContext(ctx).Model(model.Gateway{}).
		Where("gateway_id = ?", gatewayID).
		Update("status", status).
		Update("updated_at", time.Now()).
		Error
}

func (d *CoreDb) GetUserRefreshToken(ctx context.Context, userID string) (string, error) {
	var token model.Token
	err := d.db.WithContext(ctx).Model(model.Token{}).Where("user_id = ?", userID).Order("created_at desc").First(&token).Error
	return token.RefreshToken, err
}

func (d *CoreDb) AddAsyncJob(ctx context.Context, txn *gorm.DB, task *model.AsyncJob) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	if task.ExpiredAt.IsZero() {
		task.ExpiredAt = task.CreatedAt.Add(time.Hour * 2)
	}
	return txn.WithContext(ctx).Model(model.AsyncJob{}).Create(task).Error
}

func (d *CoreDb) GetAsyncJob(ctx context.Context, requestId string) (model.AsyncJob, error) {
	var ret model.AsyncJob
	err := d.db.WithContext(ctx).Model(model.AsyncJob{}).Where("request_id = ?", requestId).First(&ret).Error
	if ret.ExpiredAt.Before(time.Now()) && ret.Status != "completed" {
		ret.Status = "failed"
		// update
		d.db.WithContext(ctx).Where("request_id = ?", requestId).Save(&model.AsyncJob{
			Status: "failed",
		})
	}
	return ret, err
}

func (d *CoreDb) SetAsyncJobStatus(ctx context.Context, txn *gorm.DB, requestId string, status string) error {
	return txn.WithContext(ctx).Model(model.AsyncJob{}).Where("request_id = ?", requestId).
		Save(&model.AsyncJob{
			Status:    status,
			UpdatedAt: time.Now(),
		}).Error
}

func (d *CoreDb) DeleteGatewayHostInfo(ctx context.Context, txn *gorm.DB, id string) error {
	return txn.WithContext(ctx).Model(model.GatewayHost{}).Where("host_id = ?", id).Update("deleted_at", time.Now()).Error
}

func (d *CoreDb) GetErrorLogByErrorID(ctx context.Context, errID string) (model.ErrorLogs, error) {
	var ret model.ErrorLogs
	err := d.db.WithContext(ctx).Where("err_id = ?", errID).First(&ret).Error
	return ret, err
}

func (d *CoreDb) GetErrorLogs(ctx context.Context, gatewayid string, from, to time.Time, limit int) ([]model.ErrorLogs, error) {
	var ret []model.ErrorLogs
	var err error
	if to.IsZero() {
		// No need to filter time
		err = d.db.WithContext(ctx).Where("gateway_id = ?", gatewayid).Limit(limit).Find(&ret).Error
	} else {
		err = d.db.WithContext(ctx).Where("gateway_id = ? and created_at >= ? and created_at <= ?", gatewayid, from.Unix(), to.Unix()).Limit(limit).Find(&ret).Error
	}
	return ret, err
}

func (d *CoreDb) ListGatewayHostInfoByUserID(ctx context.Context, userID string) ([]model.GatewayHost, error) {
	var ret []model.GatewayHost
	err := d.db.WithContext(ctx).Where("user_id = ?", userID).Find(&ret).Error
	return ret, err
}

func (d *CoreDb) UpdateGatewayHostInfo(ctx context.Context, txn *gorm.DB, hostid string, info *model.GatewayHost) error {
	info.UpdatedAt = time.Now()
	return txn.WithContext(ctx).Model(model.GatewayHost{}).
		Where("host_id = ?", hostid).
		Save(info).
		Update("updated_at", time.Now()).
		Error
}

func (d *CoreDb) GetGatewayHostInfo(ctx context.Context, id string) (model.GatewayHost, error) {
	var ret model.GatewayHost
	err := d.db.WithContext(ctx).Where("host_id = ?", id).First(&ret).Error
	return ret, err
}

func (d *CoreDb) AddGatewayHostInfo(ctx context.Context, txn *gorm.DB, info *model.GatewayHost) error {
	if txn == nil {
		return errs.ErrNeedTxn
	}
	info.CreatedAt = time.Now()
	info.UpdatedAt = time.Now()
	return txn.WithContext(ctx).Create(info).Error
}

func (d *CoreDb) SaveToken(ctx context.Context, tk *model.Token) error {
	tk.CreatedAt = time.Now()
	tk.UpdatedAt = time.Now()
	return d.db.WithContext(ctx).Create(tk).Error
}
func (d *CoreDb) SaveTokenOauth2(ctx context.Context, tk *oauth2.Token, userID string) error {
	var mToken model.Token
	mToken.AccessToken = tk.AccessToken
	mToken.RefreshToken = tk.RefreshToken
	mToken.AccessTokenExpiredAt = tk.Expiry.Second()
	mToken.RefreshTokenExpiredAt = tk.Expiry.Second()
	mToken.UserId = userID
	mToken.CreatedAt = time.Now()
	mToken.UpdatedAt = time.Now()
	return d.SaveToken(ctx, &mToken)
}
func (d *CoreDb) QueryUser(ctx context.Context, userName, uuid string) (ret model.User, err error) {
	// one of these two methods will return the first user
	err = d.db.WithContext(ctx).Model(model.User{}).Where("username = ? or user_id = ?", userName, uuid).First(&ret).Error
	return
}

func (d *CoreDb) GatewayIDGetUserID(ctx context.Context, id string) (string, error) {
	var ret model.Gateway
	err := d.db.WithContext(ctx).Model(model.Gateway{}).Where("id = ?", id).First(&ret).Error
	return ret.UserId, err
}

func (d *CoreDb) AgentIDGetGatewayID(ctx context.Context, id string) (string, error) {
	var ret model.Agent
	err := d.db.WithContext(ctx).Model(model.Agent{}).Where("id = ?", id).First(&ret).Error
	return ret.GatewayId, err
}

func (d *CoreDb) AddAgent(ctx context.Context, txn *gorm.DB, gatewayID string, agent model.Agent) error {
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()
	return txn.WithContext(ctx).Create(&agent).Error
}

func (d *CoreDb) UpdateAgent(ctx context.Context, txn *gorm.DB, agent model.Agent) error {
	agent.UpdatedAt = time.Now()
	return txn.WithContext(ctx).Model(model.Agent{}).Where("agent_id = ?", agent.AgentId).Save(&agent).Error
}

func (d *CoreDb) DeleteAgent(ctx context.Context, txn *gorm.DB, id string) error {
	var _db *gorm.DB
	if txn == nil {
		_db = d.db
	} else {
		_db = txn
	}
	return _db.WithContext(ctx).
		Model(model.Agent{}).
		Where("agent_id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (d *CoreDb) GetAgentList(ctx context.Context, gatewayID string) ([]model.Agent, error) {
	var ret []model.Agent
	err := d.db.WithContext(ctx).Model(model.Agent{}).Where("gateway_id = ?", gatewayID).Find(&ret).Error
	return ret, err
}

func (d *CoreDb) GetGatewayInfo(ctx context.Context, id string) (*model.Gateway, error) {
	var ret model.Gateway
	return &ret, d.db.WithContext(ctx).Model(model.Gateway{}).Where("gateway_id = ?", id).First(&ret).Error
}

func (d *CoreDb) AddGateway(ctx context.Context, txn *gorm.DB, userID string, gateway model.Gateway) error {
	gateway.UserId = userID
	gateway.CreatedAt = time.Now()
	gateway.UpdatedAt = time.Now()
	if txn == nil {
		return errs.ErrNeedTxn
	}
	return txn.WithContext(ctx).Create(&gateway).Error
}

func (d *CoreDb) UpdateGateway(ctx context.Context, txn *gorm.DB, gateway model.Gateway) error {
	if txn == nil {
		return errs.ErrNeedTxn
	}
	var m model.Gateway
	txn.WithContext(ctx).Model(model.Gateway{}).Where("gateway_id = ?", gateway.GatewayID).First(&m)
	gateway.UpdatedAt = time.Now()
	gateway.CreatedAt = m.CreatedAt
	gateway.UserId = m.UserId
	gateway.Status = m.Status
	gateway.BindHost = m.BindHost
	gateway.ID = m.ID
	return txn.WithContext(ctx).Model(model.Gateway{}).Where("gateway_id = ?", gateway.GatewayID).Updates(&gateway).Error
}

func (d *CoreDb) DeleteGateway(ctx context.Context, txn *gorm.DB, id string) error {
	var _db *gorm.DB
	if txn == nil {
		_db = d.db
	} else {
		_db = txn
	}
	return _db.WithContext(ctx).
		Model(model.Gateway{}).
		Where("gateway_id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (d *CoreDb) GetAllGatewayInfo(ctx context.Context) ([]model.Gateway, error) {
	var ret []model.Gateway
	return ret, d.db.WithContext(ctx).Model(model.Gateway{}).Find(&ret).Error
}

func (d *CoreDb) DeleteAgentGatherData(ctx context.Context, txn *gorm.DB, id string, timeStart int64, timeEnd int64) error {
	if txn != nil {
		return txn.WithContext(ctx).Where("agent_id = ? and created_at >= ? and created_at <= ?", id, timeStart, timeEnd).Delete(&model.Data{}).Error
	}
	return d.db.WithContext(ctx).Where("agent_id = ? and created_at >= ? and created_at <= ?", id, timeStart, timeEnd).Delete(&model.Data{}).Error
}

func (d *CoreDb) Begin() *gorm.DB {
	return d.db.Begin()
}

func (d *CoreDb) Commit(txn *gorm.DB) {
	txn.Commit()
}

func (d *CoreDb) Rollback(txn *gorm.DB) {
	txn.Rollback()
}

func (d *CoreDb) GetAgentInfo(id string) (*model.Agent, error) {
	var agent model.Agent
	return &agent, d.db.Model(agent).Where("agent_id = ?", id).First(&agent).Error
}

func (d *CoreDb) Weight() uint16 {
	return 3
}

func (d *CoreDb) Version() string {
	return "dev"
}

func (d *CoreDb) IsGatewayIdExists(id string) bool {
	return d.db.Where("gateway_id = ?", id).First(&model.Gateway{}).Error == nil
}
func (d *CoreDb) StoreAgentGatherData(ctx context.Context, txn *gorm.DB, id, content string) error {
	data := &model.Data{AgentID: id, Content: content, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if txn != nil {
		return txn.WithContext(ctx).Create(data).Error
	}
	return d.db.WithContext(ctx).Create(data).Error
}
func (d *CoreDb) GetDataCleaner(id string) (*model.Clean, error) {
	var ret model.Clean
	return &ret, d.db.Where("agent_id = ?", id).First(&ret).Error
}
func (d *CoreDb) CreateTask(ctx context.Context, txn *gorm.DB, id string, task task.Task) error {
	task.TaskID = id
	return txn.WithContext(ctx).Create(&task).Error
}
func (d *CoreDb) TaskUpdateOperationType(ctx context.Context, txn *gorm.DB, id string, operationType task.Operation) error {
	return txn.WithContext(ctx).Model(&task.Task{}).Where("id = ?", id).Update("operation", operationType).Error
}
func (d *CoreDb) TaskUpdateOperationCommend(ctx context.Context, txn *gorm.DB, id string, operationCommend string) error {
	return txn.WithContext(ctx).Model(&task.Task{}).Where("id = ?", id).Update("operation_commend", operationCommend).Error
}
func (*CoreDb) Name() string {
	return "Core-database-client"
}
func (d *CoreDb) Init() error {
	d.db = DefaultCoreClient().db
	return nil
}
