package edge

import (
	"context"
	"fmt"
	"sync"

	dataService "github.com/AEnjoy/IoT-lubricant/internal/edge/grpc"
	"github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/net"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	"github.com/nats-io/nats.go"
)

type app struct {
	mq mq.Mq[[]byte]

	openapi.OpenApi

	ctrl   context.Context
	cancel context.CancelFunc

	l sync.Mutex
	// for init
	config *types.EdgeSystem

	hostAddress string // 这是容器宿主机的ip:port
}

func (a *app) Run() error {
	go compressor(a.config.Algorithm, dataSetCh, compressedChan)
	go transmitter(a.config.ReportCycle, compressedChan, triggerChan, dataChan2)

	return a.StartGather(a.ctrl)
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

func UseConfig(c *types.EdgeSystem) func(*app) error {
	return func(s *app) error {
		if c == nil {
			logger.Warnln("config is nil")
		}
		s.config = c
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
func UseGRPC(bind string) func(*app) error {
	return func(a *app) error {
		if bind == "" {
			logger.Warnln("grpc bind is empty, use default")
			bind = _default.AgentDefaultBind
		}
		go dataService.NewServer(bind)
		return nil
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
