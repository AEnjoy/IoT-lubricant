package service

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/pkg/ssh"
	errCh "github.com/AEnjoy/IoT-lubricant/pkg/error"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptionCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
)

func (s *GatewayService) AddHost(ctx context.Context, info *model.GatewayHost) error {
	txn := s.db.Begin()
	errorCh := errCh.NewErrorChan()
	defer errCh.HandleErrorCh(errorCh).ErrorWillDo(func(error) {
		s.db.Rollback(txn)
	}).SuccessWillDo(func() {
		s.db.Commit(txn)
	}).Do()
	err := s.db.AddGatewayHostInfo(ctx, txn, info)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbAddGatewayFailed,
			exception.WithMsg("Failed to add gateway information to database"),
		)
		return err
	}

	_, err = ssh.NewSSHClient(info, true)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.LinkToGatewayFailed,
			exception.WithMsg("LinkToTargetHostError:"),
		)
		errorCh.Report(err, exceptionCode.LinkToGatewayFailed, "LinkToTargetHostError:", true)
		return err
	}
	return nil
}
