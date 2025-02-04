package service

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/config"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
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
