package model

import (
	"gorm.io/gorm"
)

type GatewayDbCli interface {
	GetServerInfo() ServerInfo
	IsAgentIdExists(string) bool
	GetAllAgentId() []string
	RemoveAgent(...string) bool
}

type CoreDb struct {
	db *gorm.DB
}

func (*CoreDb) Name() string {
	return "Core-database-client"
}
func (d *CoreDb) Init() error {
	d.db = DefaultCoreClient().db
	return nil
}

type GatewayDb struct {
	db *gorm.DB
}

func (*GatewayDb) Name() string {
	return "Gateway-database-client"
}
func (d *GatewayDb) Init() error {
	d.db = Gateway(nil).db
	return nil
}
func (d *GatewayDb) GetServerInfo() (s ServerInfo) {
	d.db.First(&s)
	return
}
func (d *GatewayDb) IsAgentIdExists(id string) bool {
	return d.db.Where("id = ?", id).First(&Agent{}).Error == nil
}
func (d *GatewayDb) GetAllAgentId() (retVal []string) {
	var agents []Agent
	d.db.Find(&agents)
	for _, agent := range agents {
		retVal = append(retVal, agent.Id)
	}
	return
}
func (d *GatewayDb) RemoveAgent(id ...string) bool {
	return d.db.Where("id in (?)", id).Delete(&Agent{}).Error == nil
}
