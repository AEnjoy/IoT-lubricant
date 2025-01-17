package repo

import (
	"errors"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"gorm.io/gorm"
)

var _ GatewayDbOperator = (*GatewayDb)(nil)

type GatewayDb struct {
	db *gorm.DB
}

func (d *GatewayDb) UpdateAgentInstance(txn *gorm.DB, id string, ins model.AgentInstance) error {
	return txn.Where("agent_id = ?", id).Updates(ins).Error
}

func (d *GatewayDb) GetAgent(id string) (retVal model.Agent, err error) {
	err = d.db.Where("agent_id = ?", id).First(&retVal).Error
	return
}

func (d *GatewayDb) UpdateAgent(txn *gorm.DB, agent *model.Agent) error {
	return txn.Where("agent_id = ?", agent.AgentId).Updates(agent).Error
}

func (d *GatewayDb) Begin() *gorm.DB {
	return d.db.Begin()
}

func (d *GatewayDb) Commit(txn *gorm.DB) {
	txn.Commit()
}

func (d *GatewayDb) Rollback(txn *gorm.DB) {
	txn.Rollback()
}
func (d *GatewayDb) GetAgentInstance(txn *gorm.DB, id string) model.AgentInstance {
	var ret model.AgentInstance
	if txn == nil {
		txn = d.db
	}
	txn.Where("agent_id = ?", id).First(&ret)
	return ret
}
func (d *GatewayDb) AddAgentInstance(txn *gorm.DB, ins model.AgentInstance) error {
	return txn.Create(ins).Error
}

func (*GatewayDb) Name() string {
	return "Gateway-database-client"
}
func (d *GatewayDb) Init() error {
	d.db = NewGatewayDb(nil).db
	return nil
}
func (d *GatewayDb) GetServerInfo(_ *gorm.DB) *model.ServerInfo {
	ret := model.ServerInfo{}
	if err := d.db.First(&ret).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
	}
	return &ret
}
func (d *GatewayDb) IsAgentIdExists(_ *gorm.DB, id string) bool {
	return d.db.Where("id = ?", id).First(&model.Agent{}).Error == nil
}
func (d *GatewayDb) GetAllAgentId(_ *gorm.DB) (retVal []string) {
	var agents []model.Agent
	d.db.Find(&agents)
	for _, agent := range agents {
		retVal = append(retVal, agent.AgentId)
	}
	return
}
func (d *GatewayDb) GetAllAgents(_ *gorm.DB) (agents []model.Agent, err error) {
	err = d.db.Find(&agents).Error
	return
}

func (d *GatewayDb) RemoveAgent(txn *gorm.DB, id ...string) bool {
	return txn.Where("id in (?)", id).Delete(&model.Agent{}).Error == nil
}
func (d *GatewayDb) GetAgentReportCycle(txn *gorm.DB, id string) int {
	var agent model.Agent
	d.db.Where("id = ?", id).First(&agent)
	return agent.Cycle
}
func (d *GatewayDb) GetAgentGatherCycle(txn *gorm.DB, id string) int {
	var agent model.Agent
	d.db.Where("id = ?", id).First(&agent)
	return agent.GatherCycle
}
