package monitor

import "github.com/aenjoy/iot-lubricant/services/lubricant/services"

type Api struct {
	services.IGatewayService
	services.IAgentService
}
