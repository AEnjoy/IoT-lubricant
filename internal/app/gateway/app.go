package gateway

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/app/gateway/internal/agent"
	"github.com/AEnjoy/IoT-lubricant/internal/app/gateway/internal/async"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	def "github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/crypto"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/AEnjoy/IoT-lubricant/protobuf/meta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var gatewayId string

type app struct {
	ctrl context.Context

	repo.GatewayDbOperator
	agent agent.Apis
	task  async.Task

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
	a.agent = agent.NewAgentApis(a.GatewayDbOperator)
	a.task = async.NewAsyncTask()
	a.task.SetActor(a.handelTask)

	_ = a.agentPoolInit() //todo:handel error
	go a.agentPoolAgentRegis()
	go a.agentPoolChStartService()
	//go a.agentHandelSignal()
	for {
		time.Sleep(time.Second * 5)
	}
	var errs error
	go func() {
		if err := a.grpcTaskApp(); err != nil {
			errs = errors.Join(err)
		}
	}()
	go func() {
		if err := a.grpcDataApp(); err != nil {
			errs = errors.Join(err)
		}
	}()
	return errs // gateway <--> core
}

func SetGatewayId(id string) func(*app) error {
	return func(s *app) error {
		gatewayId = id
		s.ctrl = context.WithValue(context.Background(), types.NameGatewayID, id)
		return nil
	}
}

func UseDB(db *repo.GatewayDb) func(*app) error {
	return func(a *app) error {
		a.GatewayDbOperator = db
		return nil
	}
}

func linkToGrpcServer(address string, tls *crypto.Tls) func(*app) error {
	return func(a *app) error {
		var conn *grpc.ClientConn
		var err error
		kacp := keepalive.ClientParameters{
			Time:                10 * time.Second, // 每隔 10 秒发送一次心跳
			Timeout:             3 * time.Second,  // 心跳超时时间为 3 秒
			PermitWithoutStream: true,             // 允许在没有流的情况下发送心跳
		}
		if tls != nil && tls.Enable {
			config, err := tls.GetTLSLinkConfig()
			if err != nil {
				return err
			}
			conn, err = grpc.NewClient(address, grpc.WithTransportCredentials(config), grpc.WithKeepaliveParams(kacp))
			if err != nil {
				return err
			}
		} else {
			conn, err = grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp))
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
		if err := stream.Send(&meta.Ping{Flag: 0}); err != nil {
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

func LinkCoreServer() func(*app) error {
	return func(a *app) error {
		local := func(info *model.ServerInfo) error {
			return linkToGrpcServer(fmt.Sprintf("%s:%d", info.Host, info.Port), &info.TlsConfig)(a)
		}
		env := func(address string) error {
			return linkToGrpcServer(address, nil)(a)
		}

		address := os.Getenv(def.ENV_CORE_HOST_STR)
		info := a.GatewayDbOperator.GetServerInfo()
		if info == nil && address == "" {
			return errors.New("address should not be empty when not initialized")
		}
		if info != nil {
			if info.Host == "" || info.Port == 0 {
				logger.Error("Incorrect local configuration, starting with environment variable")
				return env(address)
			}
			logger.Info("Use local config to start")
			return local(info)
		} else {
			logger.Info("Use environment variable to start")
			return env(address)
		}
	}
}
