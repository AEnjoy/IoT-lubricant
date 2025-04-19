package response

import "github.com/aenjoy/iot-lubricant/pkg/model"

type DescriptionHostResponse struct {
	Host        model.GatewayHost `json:"host_info"`
	GatewayList []model.Gateway   `json:"gateway_list"`
}
