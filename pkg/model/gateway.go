package model

import (
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
	gatewaypb "github.com/aenjoy/iot-lubricant/protobuf/gateway"
)

type Gateway struct {
	ID          int    `json:"-" gorm:"column:id;primary_key;autoIncrement"`
	GatewayID   string `json:"gateway_id" gorm:"column:gateway_id"`
	UserId      string `json:"-" gorm:"column:user_id"` //;foreignKey:UserID
	BindHost    string `json:"_" gorm:"column:bind_host"`
	Description string `json:"description" gorm:"column:description"`

	TlsConfig string `json:"tls_config" gorm:"column:tls_config;serializer:json"`
	// host information has replaced by model.GatewayHost

	Status    string    `json:"status" gorm:"column:status;default:'created';enum('offline', 'online', 'error', 'created')"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeleteAt  time.Time `gorm:"type:datetime" json:"deleteAt"`
}

func (Gateway) TableName() string {
	return "gateway"
}
func ProxypbEditAgentRequest2Agent(pbreq *gatewaypb.EditAgentRequest) *Agent {
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
	Id        int        `json:"id" gorm:"column:id;primary_key" yaml:"-"` // uuid and token
	GatewayID string     `json:"gateway_id" gorm:"column:gateway_id" yaml:"gateway_id"`
	Host      string     `json:"host" gorm:"column:host" yaml:"host"`
	Port      int        `json:"port" gorm:"column:port" yaml:"port"`
	Tls       bool       `json:"tls" gorm:"column:tls" yaml:"tls"`
	TlsConfig crypto.Tls `json:"tls_config" gorm:"column:tls_config;type:text;serializer:json" yaml:"tlsConfig"`

	CreatedAt time.Time `json:"-" gorm:"column:created_at;type:datetime" yaml:"-"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at;type:datetime" yaml:"-"`
	DeleteAt  time.Time `gorm:"type:datetime" json:"deleteAt"`
}

func (ServerInfo) TableName() string {
	return "server_info"
}
