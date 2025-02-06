package v1

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1/gateway"
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/service"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
)

var (
	_        IGateway = (*gateway.Api)(nil)
	_gateway IGateway
)

type IGateway interface {
}

func NewGateway() IGateway {
	if _gateway == nil {
		_gateway = gateway.Api{
			DataStore:       ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore),
			IGatewayService: ioc.Controller.Get(ioc.APP_NAME_CORE_GATEWAY_SERVICE).(service.IGatewayService),
		}
	}
	return _gateway
}
