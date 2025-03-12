package services

import (
	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
	"github.com/aenjoy/iot-lubricant/services/lubricant/repo"
)

var (
	_ ioc.Object = (*AgentService)(nil)
	_ ioc.Object = (*GatewayService)(nil)
	_ ioc.Object = (*SyncTaskQueue)(nil)
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
	a.SyncTaskQueue = ioc.Controller.Get(ioc.APP_NAME_CORE_Internal_SyncTask_SERVICE).(*SyncTaskQueue)
	return nil
}

func (*AgentService) Weight() uint16 {
	return ioc.CoreGatewayAgentService
}

func (*AgentService) Version() string {
	return ""
}

func (s *SyncTaskQueue) Init() error {
	s.Mq = ioc.Controller.Get(ioc.APP_NAME_CORE_Internal_MQ_SERVICE).(mq.Mq)
	return nil
}

func (SyncTaskQueue) Weight() uint16 {
	return ioc.CoreSyncTaskSystem
}

func (SyncTaskQueue) Version() string {
	return ""
}
