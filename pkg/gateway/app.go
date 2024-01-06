package gateway

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/nats"
	"google.golang.org/grpc"
)

type app struct {
	ctrl context.Context
	mq   mq.Mq[[]byte]

	dataCli    *data
	deviceList *sync.Map

	*clientMq

	model.GatewayDbCli

	gatewayId string
	port      string
	grpc      *grpc.Server
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
		a.grpc = grpcServer
		return nil
	}
}
func SetPort(port string) func(*app) error {
	return func(s *app) error {
		s.port = port
		return nil
	}
}
func SetGatewayId(gatewayId string) func(*app) error {
	return func(s *app) error {
		s.gatewayId = gatewayId
		return nil
	}
}

func UseDB(db *model.GatewayDb) func(*app) error {
	return func(a *app) error {
		a.GatewayDbCli = db
		return nil
	}
}
