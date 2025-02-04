package service

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/crypto"
)

var _ ioc.Object = (*GatewayService)(nil)

type GatewayService struct {
	db    repo.CoreDbOperator
	store *datastore.DataStore
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

type IGatewayService interface {
	AddHost(ctx context.Context, info *model.GatewayHost) error
	EditHost(ctx context.Context, hostid string, info *model.GatewayHost) error
	GetHost(ctx context.Context, hostid string) (model.GatewayHost, error)
	UserGetHosts(ctx context.Context, userid string) ([]model.GatewayHost, error)
	DeployGatewayInstance(ctx context.Context, hostid string, description string, tls *crypto.Tls) (string, error)
}
