package request

import "github.com/aenjoy/iot-lubricant/pkg/types/crypto"

type AddGatewayHostRequest struct {
	Host        string `json:"host"` // ip:port
	Description string `json:"description"`

	UserName         string `json:"username"`
	PassWd           string `json:"password,omitempty"`
	PrivateKey       string `json:"private_key,omitempty"`
	CustomPrivateKey bool   `json:"custom_private_key,omitempty"` // 是否自定义私钥,如果真 则使用参数传递的 PrivateKey，否则使用本机 PrivateKey
}
type AddGatewayRequest struct {
	AddGatewayHostRequest
	TlsConfig *crypto.Tls `json:"tls_config"`
}
type RemoveGatewayRequest struct {
	GatewayID      string `json:"gateway_id"`
	RemoveAgents   bool   `json:"remove_agents"`
	RemoveHostInfo bool   `json:"remove_host_info"`
}
