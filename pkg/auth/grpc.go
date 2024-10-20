package auth

// grpc-server-auth
import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
)

var _ ioc.Object = (*InterceptorImpl)(nil)

type InterceptorImpl struct {
	db types.CoreDbOperator
}

func (i *InterceptorImpl) Init() error {
	cli := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(types.CoreDbOperator)
	i.db = cli
	return nil
}

func (i *InterceptorImpl) Weight() uint16 {
	return ioc.CoreGrpcAuthInterceptor
}

func (i *InterceptorImpl) Version() string {
	return "dev"
}

func NewClientPerRPCCredentials(gatewayId string) *ClientPerRPCCredentials {
	return &ClientPerRPCCredentials{
		gatewayId: gatewayId,
	}
}

type ClientPerRPCCredentials struct {
	gatewayId string
}

// GetRequestMetadata gets the current request metadata, refreshing tokens
// if required. This should be called by the transport layer on each
// request, and the data should be populated in headers or other
// context. If a status code is returned, it will be used as the status for
// the RPC (restricted to an allowable set of codes as defined by gRFC
// A54). uri is the URI of the entry point for the request.  When supported
// by the underlying implementation, ctx can be used for timeout and
// cancellation. Additionally, RequestInfo data will be available via ctx
// to this call.  TODO(zhaoq): Define the set of the qualified keys instead
// of leaving it as an arbitrary string.
func (c *ClientPerRPCCredentials) GetRequestMetadata(
	ctx context.Context,
	uri ...string,
) (
	map[string]string,
	error,
) {

	return map[string]string{
		"gateway_id": c.gatewayId,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires
// transport security.
func (c *ClientPerRPCCredentials) RequireTransportSecurity() bool {
	return false
}

// Req/Resp 拦截器
func (i *InterceptorImpl) UnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("get auth failed")
	}
	ul := md.Get("gateway_id")
	if len(ul) == 0 {
		return nil, fmt.Errorf("gateway_id not present")
	}

	if !i.db.IsGatewayIdExists(ul[0]) {
		return nil, fmt.Errorf("error gateway client")
	}

	// 响应后的处理
	return handler(ctx, req)
}

func (i *InterceptorImpl) StreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	logger.Info(metadata.FromIncomingContext(ss.Context()))
	return nil
}
