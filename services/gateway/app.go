package gateway

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	def "github.com/aenjoy/iot-lubricant/pkg/default"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/types"
	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	"github.com/aenjoy/iot-lubricant/services/gateway/repo"
	"github.com/aenjoy/iot-lubricant/services/gateway/services/agent"
	"github.com/aenjoy/iot-lubricant/services/gateway/services/async"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var gatewayId string

type app struct {
	ctrl       context.Context
	hostConfig *model.ServerInfo

	repo.IGatewayDb
	agent agent.Apis
	task  async.Task

	grpcClient corepb.CoreServiceClient //grpc
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
	a.agent = agent.NewAgentApis(a.IGatewayDb)
	a.agent.SetReporter(_reportMessage)
	a.task = async.NewAsyncTask()
	a.task.SetActor(a.handelTask)
	agent.SetErrorHandelFunc(a._handelAgentControlError)

	go a.grpcDataApp()
	go a.grpcReportApp()

	go func() {
		_ = a.grpcPingApp()
	}()

	return a.grpcTaskApp() // gateway <--> core
}

var _ctxLock sync.Mutex

func (a *app) _getRequestContext(ctx context.Context) context.Context {
	_ctxLock.Lock()
	defer _ctxLock.Unlock()

	if a.ctrl == nil {
		if ctx == nil {
			ctx = context.Background()
		}
		md := metadata.New(map[string]string{string(types.NameGatewayID): gatewayId})
		a.ctrl = metadata.NewOutgoingContext(ctx, md)
	}
	return a.ctrl
}
func SetGatewayId(id string) func(*app) error {
	return func(s *app) error {
		gatewayId = id
		s._getRequestContext(nil)
		return nil
	}
}

func UseDB(db *repo.GatewayDb) func(*app) error {
	return func(a *app) error {
		a.IGatewayDb = db
		if a.hostConfig != nil {
			txn := db.Begin()
			if err := db.AddOrUpdateServerInfo(txn, a.hostConfig); err != nil {
				return err
			}
			if err := txn.Commit().Error; err != nil {
				return err
			}
		}
		return nil
	}
}

func linkToGrpcServer(address string, tls *crypto.Tls) func(*app) error {
	return func(a *app) error {
		var conn *grpc.ClientConn
		var err error
		//kacp := keepalive.ClientParameters{
		//	Time:                30 * time.Second, // 每隔 10 秒发送一次心跳
		//	Timeout:             5 * time.Second,
		//	PermitWithoutStream: true, // 允许在没有流的情况下发送心跳
		//}
		if tls != nil && tls.Enable {
			config, err := tls.GetTLSLinkConfig()
			if err != nil {
				return err
			}

			conn, err = grpc.NewClient(address,
				grpc.WithTransportCredentials(config),
				//grpc.WithKeepaliveParams(kacp),
			)
			if err != nil {
				return err
			}
		} else {
			conn, err = grpc.NewClient(address,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				//grpc.WithKeepaliveParams(kacp),
			)
			if err != nil {
				return err
			}
		}

		a.grpcClient = corepb.NewCoreServiceClient(conn)

		// ping stream
		stream, err := a.grpcClient.Ping(a.ctrl)
		if err != nil {
			logger.Errorf("Failed to send ping request to server: %v", err)
			return err
		}

		if err := stream.Send(&metapb.Ping{Flag: 0}); err != nil {
			logger.Errorf("Failed to send ping request to server: %v", err)
			return err
		}
		resp, err := stream.Recv()
		if err != nil {
			logger.Errorf("Failed to receive response from server: %v", err)
			return err
		}
		if resp.GetFlag() != 1 {
			return errors.New("lubricant server not ready")
		}

		return stream.CloseSend()
	}
}

func LinkCoreServer() func(*app) error {
	return func(a *app) error {
		address := os.Getenv(def.ENV_CORE_HOST_STR)
		port := os.Getenv(def.ENV_CORE_PORT_STR)
		info := a.IGatewayDb.GetServerInfo(nil)
		if info == nil && (address == "" || port == "") {
			return errors.New("address should not be empty when not initialized")
		}

		local := func(info *model.ServerInfo) error {
			logger.Debugf("Use local config to start")
			return linkToGrpcServer(fmt.Sprintf("%s:%d", info.Host, info.Port), &info.TlsConfig)(a)
		}
		env := func(address, port string) error {
			portI, _ := strconv.Atoi(port)
			info = &model.ServerInfo{
				GatewayID: gatewayId,
				Host:      address,
				Port:      portI,
			}
			if port == "" {
				logger.Debugf("Use environment variable to start")
				return linkToGrpcServer(address, nil)(a)
			}
			logger.Debugf("Use environment variable to start")
			return linkToGrpcServer(fmt.Sprintf("%s:%s", address, port), nil)(a)
		}

		defer func() {
			txn := a.IGatewayDb.Begin()
			if err := a.IGatewayDb.AddOrUpdateServerInfo(txn, info); err != nil {
				logger.Error("Failed to update server info: %v", err)
				return
			}
			if err := txn.Commit().Error; err != nil {
				logger.Error("Failed to commit transaction: %v", err)
				return
			}
		}()
		if info != nil {
			if info.Host == "" || info.Port == 0 {
				logger.Error("Incorrect local configuration, starting with environment variable")
				return env(address, port)
			}
			logger.Info("Use local config to start")
			return local(info)
		} else {
			logger.Info("Use environment variable to start")
			return env(address, port)
		}
	}
}
func UseServerInfo(c *model.ServerInfo) func(*app) error {
	return func(a *app) error {
		a.hostConfig = c
		return nil
	}
}
func UseGrpcDebugServer() func(*app) error {
	return func(a *app) error {
		if os.Getenv(def.ENV_RUNNING_LEVEL) == "debug" && os.Getenv(def.ENV_ENABLE_GRPC_DEBUG_SERVER) == "true" {
			bind, ok := os.LookupEnv(def.ENV_GRPC_DEBUG_SERVER_PORT)
			if !ok {
				return fmt.Errorf("grpc debug server port not found, please set %s", def.ENV_GRPC_DEBUG_SERVER_PORT)
			}
			go NewDebugServer(fmt.Sprintf(":%s", bind))
		}
		return nil
	}
}
