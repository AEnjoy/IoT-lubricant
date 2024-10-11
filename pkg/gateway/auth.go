package gateway

// grpcServer-server-auth
import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
)

func NewClientPerRPCCredentials(agentId string) *ClientPerRPCCredentials {
	return &ClientPerRPCCredentials{
		agentId: agentId,
	}
}

type ClientPerRPCCredentials struct {
	agentId string
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
		"agent_id": c.agentId,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires
// transport security.
func (c *ClientPerRPCCredentials) RequireTransportSecurity() bool {
	return false
}
func NewInterceptorImpl() *InterceptorImpl {
	return &InterceptorImpl{}
}

type InterceptorImpl struct {
	db model.GatewayDbCli
}

// Req/Resp 拦截器
func (i *InterceptorImpl) UnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("get auth failed")
	}
	ul := md.Get("agent_id")
	if len(ul) == 0 {
		return nil, fmt.Errorf("agent_id not present")
	}

	if !i.db.IsAgentIdExists(ul[0]) {
		return nil, fmt.Errorf("error agent client")
	}

	// 响应后的处理
	return handler(ctx, req)
}

func (i *InterceptorImpl) StreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	logger.Info(metadata.FromIncomingContext(ss.Context()))
	return nil
}
