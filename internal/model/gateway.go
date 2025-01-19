package model

import (
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/types/crypto"
	"github.com/AEnjoy/IoT-lubricant/protobuf/proxy"
)

type Agent struct {
	ID          int    `json:"id" gorm:"column:id;primary_key,autoIncrement"`
	AgentId     string `json:"agent_id" gorm:"column:agent_id"` // agent id
	GatewayId   string `json:"gateway_id" gorm:"column:gateway_id"`
	Description string `json:"description" gorm:"column:description"`
	Cycle       int    `json:"cycle" gorm:"column:cycle"`               //上报周期 默认30 单位：秒
	GatherCycle int    `json:"gather_cycle" gorm:"column:gather_cycle"` //采集周期 默认1 单位：秒
	Address     string `json:"address" gorm:"column:address"`           //container IP:PORT

	Algorithm string `json:"algorithm" gorm:"column:algorithm"`
	//APIList     []DeviceAPI `json:"api_list" gorm:"column:api_list;serializer:json"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func ProxypbEditAgentRequest2Agent(pbreq *proxy.EditAgentRequest) *Agent {
	return &Agent{
		AgentId:     pbreq.GetAgentId(),
		GatewayId:   pbreq.GetInfo().GetGatewayID(),
		Description: pbreq.GetInfo().GetDescription(),
		//Cycle:        pbreq.GetInfo().GetCycle(),
		GatherCycle: int(pbreq.GetInfo().GetGatherCycle()),
		//Address:     pbreq.Address, // no need
		Algorithm: pbreq.GetInfo().GetAlgorithm(),
		UpdatedAt: time.Now(),
	}
}
func (Agent) TableName() string {
	return "agent"
}

type ServerInfo struct { // Gateway system config
	GatewayId string     `json:"id" gorm:"column:id"` // uuid and token
	Host      string     `json:"host" gorm:"column:host"`
	Port      int        `json:"port" gorm:"column:port"`
	Tls       bool       `json:"tls" gorm:"column:tls"`
	TlsConfig crypto.Tls `json:"tls_config" gorm:"column:tls_config;type:json"`

	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at"`
}

func (ServerInfo) TableName() string {
	return "server_info"
}
