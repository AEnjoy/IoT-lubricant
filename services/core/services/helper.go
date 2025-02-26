package services

import (
	"github.com/AEnjoy/IoT-lubricant/internal/pkg/ssh"
	model2 "github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptionCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	"github.com/AEnjoy/IoT-lubricant/services/core/config"
)

func (s *GatewayService) getHostInfo() *model2.ServerInfo {
	systemConfig := config.GetConfig()
	return &model2.ServerInfo{
		Host:      systemConfig.Domain,
		Port:      systemConfig.GrpcPort,
		Tls:       systemConfig.GRPCTls,
		TlsConfig: systemConfig.Tls,
	}
}
func (s *GatewayService) checkSSHLinker(info *model2.GatewayHost) error {
	_, err := ssh.NewSSHClient(info, true)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.LinkToGatewayFailed,
			exception.WithMsg("LinkToTargetHostError:"),
		)
	}
	return err
}
