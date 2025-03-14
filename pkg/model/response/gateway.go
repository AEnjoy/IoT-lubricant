package response

import "github.com/aenjoy/iot-lubricant/pkg/model"

type AddGatewayResponse struct {
	GatewayID string `json:"gateway_id"`
}
type DescriptionGatewayResponse struct {
	Gateway *model.Gateway `json:"gateway_info"`
	Agents  []model.Agent  `json:"agent_list"`
}
type ListGatewayResponse struct {
	Gateways []model.Gateway `json:"gatewayList"`
}
