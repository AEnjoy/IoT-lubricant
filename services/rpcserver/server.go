package rpcserver

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/grpc/middleware"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/aenjoy/iot-lubricant/services/corepkg/auth"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/panjf2000/ants/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type grpcServer struct {
	*grpc.Server
	PbCoreServiceImpl
}
type PbCoreServiceImpl struct {
	corepb.UnimplementedCoreServiceServer
	*datastore.DataStore
	pool              *ants.Pool
	getProjectIdMutex sync.Mutex
}

func (a *app) grpcInit() (serve func() error, err error) {
	middlewares := &auth.InterceptorImpl{Db: a.ICoreDb}
	var server *grpc.Server

	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     time.Minute * 10, // 连接空闲超过 10 分钟则关闭
		MaxConnectionAge:      time.Hour * 2,    // 连接存活超过 2 小时则关闭
		MaxConnectionAgeGrace: time.Minute * 5,  // 连接关闭前的宽限期
		Time:                  time.Minute * 1,  // 每隔 1 分钟发送一次 ping
		Timeout:               time.Second * 20, // ping 超时时间为 20 秒
	}
	kaep := keepalive.EnforcementPolicy{
		MinTime:             time.Minute * 5, // 连接建立后至少 5 分钟才能发送 ping
		PermitWithoutStream: true,            // 允许在没有流的情况下发送 ping
	}
	if a.Tls != nil && a.Tls.Enable {
		grpcTlsOption, err := a.Tls.GetServerTlsConfig()
		if err != nil {
			return nil, err
		}
		server = grpc.NewServer(
			grpcTlsOption,
			grpc.KeepaliveParams(kasp),
			grpc.KeepaliveEnforcementPolicy(kaep),
			grpc.MaxRecvMsgSize(1024*1024*100), // 100 MB
			grpc.MaxSendMsgSize(1024*1024*100), // 100 MB
			grpc.ChainStreamInterceptor(middlewares.StreamServerInterceptor),
			grpc.ChainUnaryInterceptor(
				//middlewares.UnaryServerInterceptor,
				middleware.GetLoggerInterceptor(),
				middleware.GetRecovery(middleware.GetRegistry(middleware.GetSrvMetrics())),
			),
		)
	} else {
		server = grpc.NewServer(
			grpc.KeepaliveParams(kasp),
			grpc.KeepaliveEnforcementPolicy(kaep),
			grpc.MaxRecvMsgSize(1024*1024*100), // 100 MB
			grpc.MaxSendMsgSize(1024*1024*100), // 100 MB
			grpc.ChainStreamInterceptor(middlewares.StreamServerInterceptor),
			grpc.ChainUnaryInterceptor(middlewares.UnaryServerInterceptor),
		)
	}
	corepb.RegisterCoreServiceServer(server, a)
	a.Server = server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return nil, err
	}
	return func() error {
		return server.Serve(lis)
	}, nil
}
