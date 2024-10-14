package app

import (
	"context"
	"io"

	"github.com/AEnjoy/IoT-lubricant/pkg/auth"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"google.golang.org/grpc"
)

type PbCoreServiceImpl struct {
	core.UnimplementedCoreServiceServer
}

func (PbCoreServiceImpl) Ping(grpc.BidiStreamingServer[core.Ping, core.Ping]) error {
	return nil
}
func (PbCoreServiceImpl) GetTask(grpc.BidiStreamingServer[core.Task, core.Task]) error {
	return nil
}
func (PbCoreServiceImpl) PushMessageId(context.Context, *core.MessageIdInfo) (*core.MessageIdInfo, error) {
	return nil, nil
}
func (PbCoreServiceImpl) PushData(d grpc.BidiStreamingServer[core.Data, core.Data]) error {
	for {
		data, err := d.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// 由于数据处理需要消耗一定时间，所以使用goroutine处理
		go HandelRecvData(data)

		err = d.Send(&core.Data{
			MessageId: data.MessageId,
		})
		if err != nil {
			return err
		}
	}
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
	core.RegisterCoreServiceServer(server, &PbCoreServiceImpl{})
	return server, nil
}
