package repo

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/errs"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"gorm.io/gorm"
)

var _ CoreDbOperator = (*CoreDb)(nil)

type CoreDb struct {
	db *gorm.DB
}

func (d *CoreDb) GetGatewayHostInfo(ctx context.Context, id string) (model.GatewayHost, error) {
	var ret model.GatewayHost
	var err error
	err = d.db.WithContext(ctx).Where("host_id = ?", id).First(&ret).Error
	return ret, err
}

func (d *CoreDb) AddGatewayHostInfo(ctx context.Context, txn *gorm.DB, info *model.GatewayHost) error {
	if txn == nil {
		return errs.ErrNeedTxn
	}
	return txn.WithContext(ctx).Create(info).Error
}

func (d *CoreDb) SaveToken(ctx context.Context, tk *model.Token) error {
	return d.db.WithContext(ctx).Create(tk).Error
}
func (d *CoreDb) SaveTokenOauth2(ctx context.Context, tk *oauth2.Token, userID string) error {
	var mToken model.Token
	mToken.AccessToken = tk.AccessToken
	mToken.RefreshToken = tk.RefreshToken
	mToken.AccessTokenExpiredAt = tk.Expiry.Second()
	mToken.RefreshTokenExpiredAt = tk.Expiry.Second()
	mToken.UserId = userID
	return d.SaveToken(ctx, &mToken)
}
func (d *CoreDb) QueryUser(ctx context.Context, userName, uuid string) (ret model.User, err error) {
	// one of these two methods will return the first user
	err = d.db.Where("username = ? or user_id = ?", userName, uuid).First(&ret).Error
	return
}

func (d *CoreDb) GatewayIDGetUserID(ctx context.Context, id string) (string, error) {
	var ret model.Gateway
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&ret).Error
	return ret.UserId, err
}

func (d *CoreDb) AgentIDGetGatewayID(ctx context.Context, id string) (string, error) {
	var ret model.Agent
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&ret).Error
	return ret.GatewayId, err
}

func (d *CoreDb) AddAgent(ctx context.Context, txn *gorm.DB, gatewayID string, agent model.Agent) error {
	return txn.WithContext(ctx).Create(&agent).Error
}

func (d *CoreDb) UpdateAgent(ctx context.Context, txn *gorm.DB, agent model.Agent) error {
	return txn.WithContext(ctx).Where("id = ?", agent.AgentId).Save(&agent).Error
}

func (d *CoreDb) DeleteAgent(ctx context.Context, txn *gorm.DB, id string) error {
	if txn == nil {
		return d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Agent{}).Error
	}
	return txn.WithContext(ctx).Where("id = ?", id).Delete(&model.Agent{}).Error
}

func (d *CoreDb) GetAgentList(ctx context.Context, gatewayID string) ([]model.Agent, error) {
	var ret []model.Agent
	err := d.db.WithContext(ctx).Where("gateway_id = ?", gatewayID).Find(&ret).Error
	return ret, err
}

func (d *CoreDb) GetGatewayInfo(ctx context.Context, id string) (*model.Gateway, error) {
	var ret model.Gateway
	return &ret, d.db.WithContext(ctx).Where("id = ?", id).First(&ret).Error
}

func (d *CoreDb) AddGateway(ctx context.Context, txn *gorm.DB, userID string, gateway model.Gateway) error {
	gateway.UserId = userID
	if txn == nil {
		return errs.ErrNeedTxn
	}
	return txn.WithContext(ctx).Create(&gateway).Error
}

func (d *CoreDb) UpdateGateway(ctx context.Context, txn *gorm.DB, gateway model.Gateway) error {
	if txn == nil {
		return errs.ErrNeedTxn
	}
	return txn.WithContext(ctx).Where("id = ?", gateway.GatewayID).Save(&gateway).Error
}

func (d *CoreDb) DeleteGateway(ctx context.Context, txn *gorm.DB, id string) error {
	if txn == nil {
		return d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Gateway{}).Error
	}
	return txn.WithContext(ctx).Where("id = ?", id).Delete(&model.Gateway{}).Error
}

func (d *CoreDb) GetAllGatewayInfo(ctx context.Context) ([]model.Gateway, error) {
	var ret []model.Gateway
	return ret, d.db.WithContext(ctx).Find(&ret).Error
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
	return &agent, d.db.Where("id = ?", id).First(&agent).Error
}

func (d *CoreDb) Weight() uint16 {
	return 3
}

func (d *CoreDb) Version() string {
	return "dev"
}

func (d *CoreDb) IsGatewayIdExists(id string) bool {
	return d.db.Where("id = ?", id).First(&model.Gateway{}).Error == nil
}
func (d *CoreDb) StoreAgentGatherData(ctx context.Context, txn *gorm.DB, id, content string) error {
	data := &model.Data{AgentID: id, Content: content}
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
