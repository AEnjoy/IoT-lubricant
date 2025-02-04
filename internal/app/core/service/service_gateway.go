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

type IGatewayService interface {
	AddHost(ctx context.Context, info *model.GatewayHost) error
	EditHost(ctx context.Context, hostid string, info *model.GatewayHost) error
	GetHost(ctx context.Context, hostid string) (model.GatewayHost, error)
	UserGetHosts(ctx context.Context, userid string) ([]model.GatewayHost, error)
	DeployGatewayInstance(ctx context.Context, hostid string, description string, tls *crypto.Tls) (string, error)
}
