package mq

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
)

var _ ioc.Object = (*MqService)(nil)

func (m *MqService) Init() error {
	m.DataStore = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	m.Mq = m.DataStore.Mq
	return nil
}

func (MqService) Weight() uint16 {
	return ioc.CoreMqService
}

func (MqService) Version() string {
	return "dev"
}
