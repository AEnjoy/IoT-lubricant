package gateway

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/nats"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var gatewayId string

type app struct {
	ctrl context.Context
	mq   mq.Mq[[]byte]

	dataCli    *data
	deviceList *sync.Map

	*clientMq

	model.GatewayDbCli

	port       string
	grpcServer *grpc.Server
	grpcClient core.CoreServiceClient //grpc
}

func NewApp(opts ...func(*app) error) *app {
	var app = new(app)
	for _, opt := range opts {
		if err := opt(app); err != nil {
			logger.Fatalf("Failed to apply option: %v", err)
		}
	}
	return app
}
func (a *app) Run() error {
	a.deviceList = new(sync.Map)
	a.clientMq = new(clientMq)
	a.dataCli = new(data)
	a.dataCli.Start()

	portInt, err := strconv.Atoi(a.port)
	if err != nil {
		logger.Warn("Invalid port number, using default value.")
		portInt = 4222
	}

	server, err := nats.NewNatsServer(portInt)
	if err != nil {
		return err
	}
	go server.Start()

	a.mq, err = mq.NewNatsMq[[]byte](fmt.Sprintf("nats://127.0.0.1:%s", a.port))
	if err != nil {
		return err
	}

	err = a.initClientMq()
	if err != nil {
		panic(err)
	}

	a.Start()

	select {}
}

func UseGRPC(grpcServer *grpc.Server) func(*app) error {
	return func(a *app) error {
		a.grpcServer = grpcServer
		return nil
	}
}
func SetPort(port string) func(*app) error {
	return func(s *app) error {
		s.port = port
		return nil
	}
}
func SetGatewayId(id string) func(*app) error {
	return func(s *app) error {
		gatewayId = id
		return nil
	}
}

func UseDB(db *model.GatewayDb) func(*app) error {
	return func(a *app) error {
		a.GatewayDbCli = db
		return nil
	}
}

func LinkToGrpcServer(address string, tls *model.Tls) func(*app) error {
	return func(a *app) error {
		var conn *grpc.ClientConn
		var err error
		if tls != nil && tls.Enable {
			config, err := tls.GetTLSLinkConfig()
			if err != nil {
				return err
			}
			conn, err = grpc.NewClient(address, grpc.WithTransportCredentials(config))
		}
		conn, err = grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		a.grpcClient = core.NewCoreServiceClient(conn)

		// ping stream
		stream, err := a.grpcClient.Ping(context.Background())
		if err != nil {
			return err
		}
		if err := stream.Send(&core.PingPong{Flag: 0}); err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		if resp.GetFlag() != 1 {
			return errors.New("lubricant server not ready")
		}
		return nil
	}
}
