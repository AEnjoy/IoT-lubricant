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

	model.GatewayDbOperator

	port       string
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
	a.clientMq.ctrl = a.ctrl
	a.clientMq.deviceList = a.deviceList

	err := a.initClientMq()
	if err != nil {
		panic(err)
	}

	a.Start()
	return a.grpcApp()
}
func SetPort(port string) func(*app) error {
	return func(s *app) error {
		portInt, err := strconv.Atoi(port)
		if err != nil {
			logger.Warn("Invalid port number, using default value.")
			portInt = 4222
			port = "4222"
		}

		s.port = port

		server, err := nats.NewNatsServer(portInt)
		if err != nil {
			return err
		}
		go server.Start()

		s.mq, err = mq.NewNatsMq[[]byte](fmt.Sprintf("nats://127.0.0.1:%s", s.port))
		if err != nil {
			return err
		}
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
		a.GatewayDbOperator = db
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
			if err != nil {
				return err
			}
		} else {
			conn, err = grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return err
			}
		}

		a.grpcClient = core.NewCoreServiceClient(conn)

		// ping stream
		stream, err := a.grpcClient.Ping(context.Background())
		if err != nil {
			return err
		}
		if err := stream.Send(&core.Ping{Flag: 0}); err != nil {
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
