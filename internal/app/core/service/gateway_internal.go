package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/crypto"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptionCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	taskTypes "github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/bytedance/sonic"
)

type GatewayService struct {
	db    repo.CoreDbOperator
	store *datastore.DataStore
}

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
func (s *GatewayService) RemoveGatewayInternal(ctx context.Context, gatewayid string) error {
	txn, _, commit := s.txnHelper()
	defer commit()

	taskMq := s.store.Mq
	topic := fmt.Sprintf("/monitor/%s/%s/unregister", taskTypes.TargetGateway, gatewayid)
	err := taskMq.PublishBytes(topic, []byte("unregister gateway"))
	if err != nil {
		return exception.ErrNewException(err,
			exceptionCode.MqPublishFailed,
			exception.WithMsg("Failed to send gateway unregister signal"),
		)
	}
	t, err := taskMq.Subscribe(fmt.Sprintf("%s/response", topic))
	if err != nil {
		return exception.ErrNewException(err,
			exceptionCode.MqSubscribeFailed,
			exception.WithMsg("Failed to set gateway unregister signal"),
		)
	}

	select {
	case <-ctx.Done():
		return exception.ErrNewException(ctx.Err(),
			exceptionCode.DeadLine,
			exception.WithMsg("user request cancel"),
			exception.WithMsg("database not changed"))
	case m := <-t:
		if string(m) != "ok" {
			return exception.ErrNewException(err,
				exceptionCode.RemoveGatewayFailed,
				exception.WithMsg("gateway monitor failed to unregister this gateway"),
				exception.WithMsg(fmt.Sprintf("GatewayID: %s", gatewayid)),
				exception.WithMsg(fmt.Sprintf("Message: %s", m)),
			)
		}
	case <-time.After(10 * time.Second):
		return exception.ErrNewException(err,
			exceptionCode.RemoveGatewayFailed,
			exception.WithMsg("gateway monitor failed to unregister this gateway"),
			exception.WithMsg(fmt.Sprintf("GatewayID: %s", gatewayid)),
			exception.WithMsg("timeout"),
		)
	}
	err = s.db.DeleteGateway(ctx, txn, gatewayid)
	if err != nil {
		return exception.ErrNewException(err,
			exceptionCode.RemoveGatewayFailed,
			exception.WithMsg("Failed to delete gateway information from database"),
		)
	}
	return nil
}
func (s *GatewayService) RemoveGatewayHostInternal(ctx context.Context, hostid string) error {
	panic("implement me")
}
