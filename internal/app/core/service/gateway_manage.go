package service

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptionCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
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
