package gateway

import (
	"context"
	"net"

	"github.com/aenjoy/iot-lubricant/pkg/grpc/middleware"
	"github.com/aenjoy/iot-lubricant/pkg/logger"

	gatewaypb "github.com/aenjoy/iot-lubricant/protobuf/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gatewayDebugServer struct {
	gatewaypb.UnimplementedDebugServiceServer
}

func (gatewayDebugServer) MockCoreTask(context.Context, *gatewaypb.MockCoreTaskRequest) (*gatewaypb.MockCoreTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MockCoreTask not implemented")
}
func (gatewayDebugServer) GatewayResources(context.Context, *gatewaypb.GetGatewayResourcesApiRequest) (*gatewaypb.GetGatewayResourcesApiResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GatewayResources not implemented")
}

func NewDebugServer(bind string) {
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.GetLoggerInterceptor(),
			middleware.GetRecovery(middleware.GetRegistry(middleware.GetSrvMetrics()))))
	gatewaypb.RegisterDebugServiceServer(grpcServer, &gatewayDebugServer{})
	logger.Infoln("gateway debug-grpc-server start at: ", bind)
	panic(grpcServer.Serve(lis))
}
