package v1

import (
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/services/core/api/v1/gateway"
	"github.com/AEnjoy/IoT-lubricant/services/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/services/core/services"
	"github.com/gin-gonic/gin"
)

var (
	_        IGateway = (*gateway.Api)(nil)
	_gateway IGateway
)

type IGateway interface {
	AddHost(c *gin.Context)

	AddGatewayInternal(c *gin.Context)
	RemoveGatewayInternal(c *gin.Context)

	AgentPushTask(c *gin.Context)
	AddAgentInternal(c *gin.Context)
}

func NewGateway() IGateway {
	if _gateway == nil {
		_gateway = gateway.Api{
			DataStore:       ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore),
			IGatewayService: ioc.Controller.Get(ioc.APP_NAME_CORE_GATEWAY_SERVICE).(services.IGatewayService),
			IAgentService:   ioc.Controller.Get(ioc.APP_NAME_CORE_GATEWAY_AGENT_SERVICE).(services.IAgentService),
		}
	}
	return _gateway
}
