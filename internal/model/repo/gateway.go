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

func (d *GatewayDb) GetAgentInstance(id string) model.AgentInstance {
	var ret model.AgentInstance
	d.db.Where("agent_id = ?", id).First(&ret)
	return ret
}

func (*GatewayDb) Name() string {
	return "Gateway-database-client"
}
func (d *GatewayDb) Init() error {
	d.db = NewGatewayDb(nil).db
	return nil
}
func (d *GatewayDb) GetServerInfo() *model.ServerInfo {
	ret := model.ServerInfo{}
	if err := d.db.First(&ret).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
	}
	return &ret
}
func (d *GatewayDb) IsAgentIdExists(id string) bool {
	return d.db.Where("id = ?", id).First(&model.Agent{}).Error == nil
}
func (d *GatewayDb) GetAllAgentId() (retVal []string) {
	var agents []model.Agent
	d.db.Find(&agents)
	for _, agent := range agents {
		retVal = append(retVal, agent.AgentId)
	}
	return
}
func (d *GatewayDb) GetAllAgents() (agents []model.Agent, err error) {
	err = d.db.Find(&agents).Error
	return
}

func (d *GatewayDb) RemoveAgent(id ...string) bool {
	return d.db.Where("id in (?)", id).Delete(&model.Agent{}).Error == nil
}
func (d *GatewayDb) GetAgentReportCycle(id string) int {
	var agent model.Agent
	d.db.Where("id = ?", id).First(&agent)
	return agent.Cycle
}
func (d *GatewayDb) GetAgentGatherCycle(id string) int {
	var agent model.Agent
	d.db.Where("id = ?", id).First(&agent)
	return agent.GatherCycle
}
