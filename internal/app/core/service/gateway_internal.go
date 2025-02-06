package service

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/crypto"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptionCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	"github.com/bytedance/sonic"
)

func (s *GatewayService) AddHostInternal(ctx context.Context, info *model.GatewayHost) error {
	txn, _, commit := s.txnHelper()
	defer commit()

	err := s.db.AddGatewayHostInfo(ctx, txn, info)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbAddGatewayFailed,
			exception.WithMsg("Failed to add gateway information to database"),
		)
		return err
	}
	return nil
}
func (s *GatewayService) AddGatewayInternal(ctx context.Context, userid, gatewayid, description string, tls *crypto.Tls) error {
	txn, errorCh, commit := s.txnHelper()
	defer commit()

	err := s.db.AddGateway(ctx, txn, userid, model.Gateway{
		GatewayID:   gatewayid,
		Description: description,
		TlsConfig: func() string {
			if tls == nil {
				return ""
			}
			marshalString, err := sonic.MarshalString(tls)
			if err != nil {
				return ""
			}
			return marshalString
		}(),
	})
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbAddGatewayFailed,
			exception.WithMsg("Failed to add gateway information to database"),
		)
		errorCh.Report(err, exceptionCode.DbAddGatewayFailed, "Failed to add gateway information to database", true)
		return err
	}
	return nil
}
