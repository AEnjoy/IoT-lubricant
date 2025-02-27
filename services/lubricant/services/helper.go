package services

import (
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/ssh"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/config"
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
func (s *GatewayService) checkSSHLinker(info *model.GatewayHost) error {
	_, err := ssh.NewSSHClient(info, true)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.LinkToGatewayFailed,
			exception.WithMsg("LinkToTargetHostError:"),
		)
	}
	return err
}
