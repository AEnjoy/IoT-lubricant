package api

import (
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/agent"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/gateway"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/log"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/monitor"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/project"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/task"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/user"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
	"github.com/aenjoy/iot-lubricant/services/corepkg/repo"

	"github.com/aenjoy/iot-lubricant/services/apiserver/services"
)

var (
	_ IGateway = (*gateway.Api)(nil)
	_ IUser    = (*user.Api)(nil)
	_ IMonitor = (*monitor.Api)(nil)
	_ IAgent   = (*agent.Api)(nil)
	_ ITask    = (*task.Api)(nil)
	_ ILog     = (*log.Api)(nil)
	_ IProject = (*project.Api)(nil)

	_gateway IGateway
	_user    IUser
	_auth    IAuth
	_monitor IMonitor
	_agent   IAgent
	_task    ITask
	_log     ILog
	_project IProject
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
func NewAgent() IAgent {
	if _agent == nil {
		_agent = agent.Api{
			DataStore:     ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore),
			IAgentService: ioc.Controller.Get(ioc.APP_NAME_CORE_GATEWAY_AGENT_SERVICE).(services.IAgentService),
		}
	}
	return _agent
}
func NewTask() ITask {
	if _task == nil {
		_task = task.Api{
			DataStore: ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore),
		}
	}
	return _task
}
func NewLog() ILog {
	if _log == nil {
		_log = log.Api{
			DataStore: ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore),
		}
	}
	return _log
}
func NewProject() IProject {
	if _project == nil {
		_project = project.Api{
			DataStore:       ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore),
			IProjectService: ioc.Controller.Get(ioc.APP_NAME_CORE_PROJECT_SERVICE).(services.IProjectService),
		}
	}
	return _project
}
