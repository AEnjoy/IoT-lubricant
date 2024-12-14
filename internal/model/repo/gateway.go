package repo

import (
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"gorm.io/gorm"
)

var _ GatewayDbOperator = (*GatewayDb)(nil)

type GatewayDb struct {
	db *gorm.DB
}

func (*GatewayDb) Name() string {
	return "Gateway-database-client"
}
func (d *GatewayDb) Init() error {
	d.db = NewGatewayDb(nil).db
	return nil
}
func (d *GatewayDb) GetServerInfo() (s model.ServerInfo) {
	d.db.First(&s)
	return
}
func (d *GatewayDb) IsAgentIdExists(id string) bool {
	return d.db.Where("id = ?", id).First(&model.Agent{}).Error == nil
}
func (d *GatewayDb) GetAllAgentId() (retVal []string) {
	var agents []model.Agent
	d.db.Find(&agents)
	for _, agent := range agents {
		retVal = append(retVal, agent.Id)
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
