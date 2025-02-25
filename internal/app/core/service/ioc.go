package service

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
)

var (
	_ ioc.Object = (*AgentService)(nil)
	_ ioc.Object = (*GatewayService)(nil)
)

func (s *GatewayService) Init() error {
	s.db = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(*repo.CoreDb)
	s.store = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
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
	return nil
}

func (*AgentService) Weight() uint16 {
	return ioc.CoreGatewayAgentService
}

func (*AgentService) Version() string {
	return ""
}
