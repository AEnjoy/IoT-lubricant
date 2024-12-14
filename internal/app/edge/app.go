package edge

import (
	"context"
	"fmt"

	dataService "github.com/AEnjoy/IoT-lubricant/internal/app/edge/grpc"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	def "github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/AEnjoy/IoT-lubricant/pkg/edge/config"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/compress"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/net"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	"github.com/nats-io/nats.go"
)

type app struct {
	ctrl   context.Context
	cancel context.CancelFunc

	// for init
	config *model.EdgeSystem

	hostAddress string // 这是容器宿主机的ip:port
}

func (a *app) Run() error {
	_compressor, _ = compress.NewCompressor(a.config.Algorithm)
	go DataHandler()

	//if edge.CheckConfigInvalid(a.OpenApi) {
	//	//config.GatherSignal <- a.ctrl
	//}
	return a.handelGatherSignalCh()
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

func UseConfig(c *model.EdgeSystem) func(*app) error {
	return func(s *app) error {
		if c == nil {
			logger.Warnln("config is nil")
			c = config.NullConfig()
		}
		s.config = c
		config.Config = c
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
			bind = def.AgentDefaultBind
		}
		go dataService.NewServer(bind)
		return nil
	}
}
func UseOpenApi(api openapi.OpenApi, err error) func(*app) error {
	return func(a *app) error {
		config.Config.Config = api
		if api != nil {
			config.OriginConfig = api.(*openapi.ApiInfo).OpenAPICli
		}
		_ = config.RefreshSlot()
		return err
	}
}
func UseCtrl(ctx context.Context) func(*app) error {
	return func(a *app) error {
		a.ctrl = ctx
		return nil
	}
}
func UseSignalHandler(handler func()) func(*app) error {
	return func(a *app) error {
		go handler()
		return nil
	}
}
