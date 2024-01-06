package model

import "time"

type Agent struct {
	Id          string `json:"id" gorm:"column:id;primary_key"` // agent id
	UserId      string `json:"user_id" gorm:"column:user_id"`
	Description string `json:"description" gorm:"column:description"`
	Cycle       int    `json:"cycle" gorm:"column:cycle"` //上报周期 默认30 单位：秒
	//APIList     []DeviceAPI `json:"api_list" gorm:"column:api_list;serializer:json"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Agent) TableName() string {
	return "agent"
}

type ServerInfo struct {
	UserId string `json:"user_id" gorm:"column:user_id"` // uuid and token
	Host   string `json:"host" gorm:"column:host"`
	Port   int    `json:"port" gorm:"column:port"`
	Tls    bool   `json:"tls" gorm:"column:tls"`

	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at"`
}

func (ServerInfo) TableName() string {
	return "server_info"
}
