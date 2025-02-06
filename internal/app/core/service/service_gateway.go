package service

import (
	"context"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/crypto"
	"google.golang.org/genproto/googleapis/rpc/status"
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

	GetRegisterStatus(_ context.Context, gatewayid string) *status.Status
	GetStatus(_ context.Context, gatewayid string) *status.Status

	GetErrorLogs(ctx context.Context, gatewayid string, from, to time.Time, limit int) ([]model.ErrorLogs, error)
	DescriptionError(ctx context.Context, errorID string) (model.ErrorLogs, error)

	HostGetGatewayDeployConfig(ctx context.Context, hostid string) (*model.ServerInfo, error)
	GatewayGetGatewayDeployConfig(ctx context.Context, gatewayid string) (*model.ServerInfo, error)

	HostSetGatewayDeployConfig(ctx context.Context, hostid string, info *model.ServerInfo) error
	GatewaySetGatewayDeployConfig(ctx context.Context, gatewayid string, info *model.ServerInfo) error

	// internal
	AddHostInternal(ctx context.Context, info *model.GatewayHost) error
	AddGatewayInternal(ctx context.Context, userID, gatewayID, description string, tls *crypto.Tls) error
}
