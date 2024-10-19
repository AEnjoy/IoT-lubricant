package model

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"gorm.io/gorm"
)

var _ ioc.Object = (*CoreDb)(nil)

type CoreDb struct {
	db *gorm.DB
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
func (d *CoreDb) StoreAgentGatherData(id, content string) error {
	data := &Data{AgentID: id, Content: content}
	return d.db.Model(data).Save(data).Error
}
func (d *CoreDb) GetDataCleaner(id string) (*Clean, error) {
	var ret Clean
	return &ret, d.db.Model(ret).Where("agent_id = ?", id).First(&ret).Error
}
func (*CoreDb) Name() string {
	return "Core-database-client"
}
func (d *CoreDb) Init() error {
	d.db = DefaultCoreClient().db
	return nil
}
