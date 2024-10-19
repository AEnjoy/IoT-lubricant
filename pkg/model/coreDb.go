package model

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"gorm.io/gorm"
)

var _ ioc.Object = (*CoreDb)(nil)

type CoreDb struct {
	db *gorm.DB
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
func (*CoreDb) Name() string {
	return "Core-database-client"
}
func (d *CoreDb) Init() error {
	d.db = DefaultCoreClient().db
	return nil
}
