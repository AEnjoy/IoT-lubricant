package types

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"gorm.io/gorm"
)

var _ ioc.Object = (*CoreDb)(nil)

type CoreDb struct {
	db *gorm.DB
}

func (d *CoreDb) GatewayIDGetUserID(ctx context.Context, id string) (string, error) {
	var ret Gateway
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&ret).Error
	return ret.UserId, err
}

func (d *CoreDb) AgentIDGetGatewayID(ctx context.Context, id string) (string, error) {
	var ret Agent
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&ret).Error
	return ret.GatewayId, err
}

func (d *CoreDb) AddAgent(ctx context.Context, txn *gorm.DB, gatewayID string, agent Agent) error {
	return txn.WithContext(ctx).Create(&agent).Error
}

func (d *CoreDb) UpdateAgent(ctx context.Context, txn *gorm.DB, agent Agent) error {
	return txn.WithContext(ctx).Where("id = ?", agent.Id).Save(&agent).Error
}

func (d *CoreDb) DeleteAgent(ctx context.Context, txn *gorm.DB, id string) error {
	if txn == nil {
		return d.db.WithContext(ctx).Where("id = ?", id).Delete(&Agent{}).Error
	}
	return txn.WithContext(ctx).Where("id = ?", id).Delete(&Agent{}).Error
}

func (d *CoreDb) GetAgentList(ctx context.Context, gatewayID string) ([]Agent, error) {
	var ret []Agent
	err := d.db.WithContext(ctx).Where("gateway_id = ?", gatewayID).Find(&ret).Error
	return ret, err
}

func (d *CoreDb) GetGatewayInfo(ctx context.Context, id string) (*Gateway, error) {
	var ret Gateway
	return &ret, d.db.WithContext(ctx).Where("id = ?", id).First(&ret).Error
}

func (d *CoreDb) AddGateway(ctx context.Context, txn *gorm.DB, userID string, gateway Gateway) error {
	gateway.UserId = userID
	if txn == nil {
		return ErrNeedTxn
	}
	return txn.WithContext(ctx).Create(&gateway).Error
}

func (d *CoreDb) UpdateGateway(ctx context.Context, txn *gorm.DB, gateway Gateway) error {
	if txn == nil {
		return ErrNeedTxn
	}
	return txn.WithContext(ctx).Where("id = ?", gateway.GatewayID).Save(&gateway).Error
}

func (d *CoreDb) DeleteGateway(ctx context.Context, txn *gorm.DB, id string) error {
	if txn == nil {
		return d.db.WithContext(ctx).Where("id = ?", id).Delete(&Gateway{}).Error
	}
	return txn.WithContext(ctx).Where("id = ?", id).Delete(&Gateway{}).Error
}

func (d *CoreDb) GetAllGatewayInfo(ctx context.Context) ([]Gateway, error) {
	var ret []Gateway
	return ret, d.db.WithContext(ctx).Find(&ret).Error
}

func (d *CoreDb) DeleteAgentGatherData(ctx context.Context, txn *gorm.DB, id string, timeStart int64, timeEnd int64) error {
	if txn != nil {
		return txn.WithContext(ctx).Where("agent_id = ? and created_at >= ? and created_at <= ?", id, timeStart, timeEnd).Delete(&Data{}).Error
	}
	return d.db.WithContext(ctx).Where("agent_id = ? and created_at >= ? and created_at <= ?", id, timeStart, timeEnd).Delete(&Data{}).Error
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

func (d *CoreDb) GetAgentInfo(id string) (*Agent, error) {
	var agent Agent
	return &agent, d.db.Where("id = ?", id).First(&agent).Error
}

func (d *CoreDb) Weight() uint16 {
	return ioc.CoreDB
}

func (d *CoreDb) Version() string {
	return "dev"
}

func (d *CoreDb) IsGatewayIdExists(id string) bool {
	return d.db.Where("id = ?", id).First(&Gateway{}).Error == nil
}
func (d *CoreDb) StoreAgentGatherData(ctx context.Context, txn *gorm.DB, id, content string) error {
	data := &Data{AgentID: id, Content: content}
	if txn != nil {
		return txn.WithContext(ctx).Create(data).Error
	}
	return d.db.WithContext(ctx).Create(data).Error
}
func (d *CoreDb) GetDataCleaner(id string) (*Clean, error) {
	var ret Clean
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
