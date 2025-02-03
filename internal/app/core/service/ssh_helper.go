package service

import (
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/pkg/ssh"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptionCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
)

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
