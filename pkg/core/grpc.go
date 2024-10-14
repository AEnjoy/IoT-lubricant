package core

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/pkg/auth"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PbCoreServiceImpl struct {
	gateway.UnimplementedGatewayServiceServer
}

func (PbCoreServiceImpl) Data(grpc.BidiStreamingServer[gateway.DataMessage, gateway.DataMessage]) error {
	return status.Errorf(codes.Unimplemented, "method Data not implemented")
}
func (PbCoreServiceImpl) Ping(context.Context, *gateway.PingPong) (*gateway.PingPong, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (PbCoreServiceImpl) PushMessageId(context.Context, *gateway.AgentMessageIdInfo) (*gateway.AgentMessageIdInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushMessageId not implemented")
}
func NewGrpcServer(port string, tls *model.Tls) (*grpc.Server, error) {
	middlewares := auth.NewInterceptorImpl()
	var server *grpc.Server
	if tls != nil {
		serverOption, err := tls.GetServerTlsConfig()
		if err != nil {
			return nil, err
		}
		server = grpc.NewServer(
			serverOption,
			grpc.ChainStreamInterceptor(middlewares.StreamServerInterceptor),
			grpc.ChainUnaryInterceptor(middlewares.UnaryServerInterceptor),
		)
	} else {
		server = grpc.NewServer(
			grpc.ChainStreamInterceptor(middlewares.StreamServerInterceptor),
			grpc.ChainUnaryInterceptor(middlewares.UnaryServerInterceptor),
		)
	}
	gateway.RegisterGatewayServiceServer(server, &PbCoreServiceImpl{})
	return server, nil
}
