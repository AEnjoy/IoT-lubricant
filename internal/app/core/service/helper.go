package service

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/config"
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
)

func (s *GatewayService) getHostInfo() *model.ServerInfo {
	systemConfig := config.GetConfig()
	return &model.ServerInfo{
		Host:      systemConfig.Domain,
		Port:      systemConfig.GrpcPort,
		Tls:       systemConfig.GRPCTls,
		TlsConfig: systemConfig.Tls,
	}
}
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
