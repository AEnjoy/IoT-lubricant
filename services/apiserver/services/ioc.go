package services

import (
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
	"github.com/aenjoy/iot-lubricant/services/corepkg/repo"
	"github.com/aenjoy/iot-lubricant/services/corepkg/syncQueue"
)

var (
	_ ioc.Object = (*AgentService)(nil)
	_ ioc.Object = (*GatewayService)(nil)
	_ ioc.Object = (*ProjectService)(nil)
)

func (s *GatewayService) Init() error {
	s.db = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(*repo.CoreDb)
	s.store = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	s.SyncTaskQueue = ioc.Controller.Get(ioc.APP_NAME_CORE_Internal_SyncTask_SERVICE).(*syncQueue.SyncTaskQueue)
	return nil
}

func (s *GatewayService) Weight() uint16 {
	return ioc.CoreGatewayService
}

func (s *GatewayService) Version() string {
	return ""
}

func (a *AgentService) Init() error {
	a.db = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(*repo.CoreDb)
	a.store = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	a.SyncTaskQueue = ioc.Controller.Get(ioc.APP_NAME_CORE_Internal_SyncTask_SERVICE).(*syncQueue.SyncTaskQueue)
	return nil
}

func (*AgentService) Weight() uint16 {
	return ioc.CoreGatewayAgentService
}

func (*AgentService) Version() string {
	return ""
}

func (p *ProjectService) Init() error {
	p.DataStore = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	return nil
}

func (ProjectService) Weight() uint16 {
	return ioc.CoreProjectService
}

func (ProjectService) Version() string {
	return "dev"
}
