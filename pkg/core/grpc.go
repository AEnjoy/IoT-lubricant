package core

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/AEnjoy/IoT-lubricant/pkg/auth"
	"github.com/AEnjoy/IoT-lubricant/pkg/core/config"
	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"google.golang.org/grpc"
)

var _ ioc.Object = (*Grpc)(nil)

// Grpc grpc server object client
type Grpc struct {
	*grpc.Server
	PbCoreServiceImpl
}

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

func (g *Grpc) Init() error {
	c := ioc.Controller.Get(config.APP_NAME).(*config.Config)
	middlewares := ioc.Controller.Get(ioc.APP_NAME_CORE_GRPC_AUTH_INTERCEPTOR).(*auth.InterceptorImpl)
	var server *grpc.Server

	if c.Tls.Enable {
		serverOption, err := c.Tls.GetServerTlsConfig()
		if err != nil {
			return err
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
	core.RegisterCoreServiceServer(server, g)
	g.Server = server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.GrpcPort))
	if err != nil {
		return err
	}
	go func() {
		_ = server.Serve(lis)
	}()
	return nil
}

func (Grpc) Weight() uint16 {
	return ioc.CoreGrpcServer
}

func (Grpc) Version() string {
	return "dev"
}
