package edge

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/net"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

type app struct {
	grpcClient gateway.GatewayServiceClient
	mq         mq.Mq[[]byte]
	*clientMq

	openapi.OpenApi

	ctrl   context.Context
	cancel context.CancelFunc

	errPanic chan error
	l        sync.Mutex
	// for init
	config      *model.EdgeSystem
	grpcConn    *grpc.ClientConn
	hostAddress string // 这是容器宿主机的ip:port
}

func (a *app) Run() error {
	a.errPanic = make(chan error)
	a.clientMq = new(clientMq)
	//if a.grpcClient != nil {
	//	a.grpcClient = gateway.NewGatewayServiceClient(a.grpcConn)
	//}

	go a.StartGather(a.ctrl)
	go compressor(a.config.Algorithm, dataSetCh, compressedChan)
	go transmitter(a.config.ReportCycle, compressedChan, triggerChan, dataChan2)
	//go a.clientGrpc() //grpc
	for err := range a.errPanic {
		return err
	}
	return nil
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

func UseConfig(config *model.EdgeSystem) func(*app) error {
	return func(s *app) error {
		if config == nil {
			logger.Errorln("config is nil")
			return fmt.Errorf("config is nil")
		}
		s.config = config
		return nil
	}
}

func UseHostAddress(address string) func(*app) error {
	return func(s *app) error {
		if address == "" || address == "auto" {
			gateway, err := net.GetGateway()
			if err != nil {
				logger.Errorln("Failed to get gateway: %v Use default", err)
			}
			s.hostAddress = fmt.Sprintf("%s:%d", gateway, nats.DefaultPort)
		}
		s.hostAddress = address
		return nil
	}
}
func UseGRPC(grpcServer *grpc.ClientConn, err error) func(*app) error {
	return func(a *app) error {
		if grpcServer != nil {
			return errors.New("grpcServer is nil")
		}
		a.grpcConn = grpcServer
		return err
	}
}
func UseOpenApi(api openapi.OpenApi, err error) func(*app) error {
	return func(a *app) error {
		a.OpenApi = api
		return err
	}
}
func UseCtrl(ctx context.Context) func(*app) error {
	return func(a *app) error {
		a.ctrl = ctx
		return nil
	}
}
func UseMq(mq mq.Mq[[]byte], err error) func(*app) error {
	return func(a *app) error {
		a.mq = mq
		return err
	}
}
