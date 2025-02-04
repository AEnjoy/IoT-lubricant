package service

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/pkg/ssh"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/crypto"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptionCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

func (s *GatewayService) AddHost(ctx context.Context, info *model.GatewayHost) error {
	txn, errorCh, commit := s.txnHelper()
	defer commit()

	err := s.db.AddGatewayHostInfo(ctx, txn, info)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbAddGatewayFailed,
			exception.WithMsg("Failed to add gateway information to database"),
		)
		return err
	}

	err = s.checkSSHLinker(info)
	if err != nil {
		errorCh.Report(err, exceptionCode.LinkToGatewayFailed, "LinkToTargetHostError:", true)
		return err
	}
	return nil
}
func (s *GatewayService) EditHost(ctx context.Context, hostid string, info *model.GatewayHost) error {
	txn, errorCh, commit := s.txnHelper()
	defer commit()

	hostInfo, err := s.db.GetGatewayHostInfo(ctx, hostid)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbGetGatewayFailed,
			exception.WithMsg("Failed to get gateway information from database"),
			exception.WithMsg("cannot compare gateway information"),
		)
		errorCh.Report(err, exceptionCode.DbGetGatewayFailed, "Failed to get gateway information from database", true)
		return err
	}
	if info.UserName != "" {
		hostInfo.UserName = info.UserName
	}
	if info.Host != "" {
		hostInfo.Host = info.Host
	}
	if info.PassWd != "" {
		hostInfo.PassWd = info.PassWd
	}
	if info.PrivateKey != "" {
		hostInfo.PrivateKey = info.PrivateKey
	}
	if info.Description != "" {
		hostInfo.Description = info.Description
	}
	err = s.db.UpdateGatewayHostInfo(ctx, txn, hostid, &hostInfo)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbUpdateGatewayInfoFailed,
			exception.WithMsg("Failed to update gateway information to database"),
		)
		errorCh.Report(err, exceptionCode.DbUpdateGatewayInfoFailed, "Failed to update gateway information to database", true)
		return err
	}
	return nil
}
func (s *GatewayService) GetHost(ctx context.Context, hostid string) (model.GatewayHost, error) {
	return s.db.GetGatewayHostInfo(ctx, hostid)
}
func (s *GatewayService) UserGetHosts(ctx context.Context, userid string) ([]model.GatewayHost, error) {
	return s.db.ListGatewayHostInfoByUserID(ctx, userid)
}

// DeployGatewayInstance 部署网关实例，返回gatewayID,error
func (s *GatewayService) DeployGatewayInstance(ctx context.Context,
	hostid, description string, tls *crypto.Tls) (string, error) {
	txn, errorCh, commit := s.txnHelper()
	defer commit()

	hostInfo, err := s.db.GetGatewayHostInfo(ctx, hostid)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbGetGatewayFailed,
			exception.WithMsg("Failed to get gateway information from database"),
		)
		errorCh.Report(err, exceptionCode.DbGetGatewayFailed, "Failed to get gateway information from database", true)
		return "", err
	}

	host, err := ssh.NewSSHClient(&hostInfo, false)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.LinkToGatewayFailed,
			exception.WithMsg("LinkToTargetHostError:"),
		)
		errorCh.Report(err, exceptionCode.LinkToGatewayFailed, "LinkToTargetHostError:", true)
		return "", err
	}

	gatewayID := uuid.NewString()
	serverInfo := s.getHostInfo()
	serverInfo.GatewayId = gatewayID

	err = host.DeployGateway(serverInfo)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.ErrorDeployGatewayFailed,
			exception.WithMsg("DeployGatewayFailed"),
		)
		errorCh.Report(err, exceptionCode.ErrorDeployGatewayFailed, "DeployGatewayFailed", true)
		return "", err
	}
	// todo:check gateway status

	err = s.db.AddGateway(ctx, txn, hostInfo.UserID, model.Gateway{
		GatewayID:   gatewayID,
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
		return "", err
	}
	return gatewayID, nil
}
