package api

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/gateway"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/monitor"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/user"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
	"github.com/aenjoy/iot-lubricant/services/lubricant/repo"
	"github.com/aenjoy/iot-lubricant/services/lubricant/services"
)

var (
	_ IGateway = (*gateway.Api)(nil)
	_ IUser    = (*user.Api)(nil)
	_ IMonitor = (*monitor.Api)(nil)

	_gateway IGateway
	_user    IUser
	_auth    IAuth
	_monitor IMonitor
)

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
func NewUser() IUser {
	if _user == nil {
		_user = &user.Api{Db: ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(repo.ICoreDb)}
	}
	return _user
}

func NewAuth() IAuth {
	if _auth == nil {
		_auth = &v1.Auth{Db: ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(repo.ICoreDb)}
	}
	return _auth
}
func NewMonitor() IMonitor {
	if _monitor == nil {
		_monitor = monitor.Api{
			DataStore:       ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore),
			IGatewayService: ioc.Controller.Get(ioc.APP_NAME_CORE_GATEWAY_SERVICE).(services.IGatewayService),
			IAgentService:   ioc.Controller.Get(ioc.APP_NAME_CORE_GATEWAY_AGENT_SERVICE).(services.IAgentService),
		}
	}
	return _monitor
}
