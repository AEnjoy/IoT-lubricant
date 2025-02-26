package task

import (
	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
)

type Way uint8

const (
	WayGatewayHostInfo Way = iota
	WayPreCoreInfo
)

type AddGatewayRequest struct {
	Way    Way `json:"way"`
	Config any `json:"config"` // HostInfo or PreCoreInfo
}

// 网关设备连接信息
type HostInfo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"` //todo: 密码加密
}

// 预置Core连接信息
type PreCoreInfo struct {
	GatewayID string     `json:"gateway_id"`
	Host      string     `json:"host"`
	Port      int        `json:"port"`
	Tls       bool       `json:"tls"`
	TlsConfig crypto.Tls `json:"tls_config"`
}

type RemoveGatewayRequest struct {
	ID          string `json:"id"`
	RemoveAgent bool   `json:"remove_agent"`
}
