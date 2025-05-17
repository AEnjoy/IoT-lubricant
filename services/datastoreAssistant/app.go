package datastoreAssistant

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/grpc/middleware"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/datastoreAssistant/api"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"

	"github.com/panjf2000/ants/v2"
	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type app struct {
	subscribed  map[string]struct{} //projectId-subscribed
	subMapMutex sync.Mutex

	*datastore.DataStore
	cli                  *etcd.Client
	internalThreadNumber int
	pool                 *ants.PoolWithFunc
	Tls                  *crypto.Tls

	StoringMutex sync.Mutex // 只要有任意的线程在存储数据，则不允许未消费完就退出

	svcpb.UnimplementedDataStoreServiceServer
	debug svcpb.UnimplementedDataStoreDebugServiceServer
	*grpc.Server

	Ctx    context.Context
	Cancel context.CancelFunc
}

func (a *app) Run() error {
	a.subscribed = make(map[string]struct{})
	regPool, err := ants.NewPoolWithFunc(a.internalThreadNumber, a.handelProjectIDStr, ants.WithPreAlloc(true))
	if err != nil {
		return err
	}
	a.pool = regPool

	hostname, _ := os.Hostname()
	consumerID := fmt.Sprintf("%s-%d", hostname, time.Now().UnixNano())
	projectIdBytesCh, err := a.V2mq.Subscribe(constant.DATASTORE_PROJECT)
	if err != nil {
		return err
	}

	id, err := a.registerConsumer(consumerID)
	if err != nil {
		return fmt.Errorf("[%s] Failed to register consumer: %v", consumerID, err)
	}
	logg.L.Debugf("Consumer:[%s] LeaseID: %d", consumerID, id)

	go a.watchAssignments(consumerID)
	go a.campaignForLeadership(consumerID)

	printAssignedOnceMap := make(map[int]struct{})
	_printAssignedOnceMapMutex := sync.Mutex{}
	var printAssigned = func(partitionID int) {
		_printAssignedOnceMapMutex.Lock()
		defer _printAssignedOnceMapMutex.Unlock()

		if _, ok := printAssignedOnceMap[partitionID]; !ok {
			printAssignedOnceMap[partitionID] = struct{}{}
			logg.L.Infof("[%s] Assigned partition: %d", consumerID, partitionID)
		}
	}

	for {
		select {
		case <-a.Ctx.Done():
			time.Sleep(3 * time.Second)
			os.Exit(0)
		case projectIdBytes := <-projectIdBytesCh:
			projectId := string(projectIdBytes.([]byte))
			partitionID := api.GetPartition(projectId, 64)
			if isPartitionAssigned(partitionID) {
				printAssigned(partitionID)
				a.subMapMutex.Lock()
				_, ok := a.subscribed[projectId]
				if !ok {
					err := a.pool.Invoke(projectId)
					if err != nil {
						logg.L.Errorf("[%s] Failed to invoke pool: %v", projectId, err)
					} else {
						a.subscribed[projectId] = struct{}{}
					}
				}
				a.subMapMutex.Unlock()
			}
		}
	}
}

func ExitHandel() {
	// todo: 优雅退出
}

func NewApp(opts ...func(*app) error) *app {
	var server = new(app)
	for _, opt := range opts {
		if err := opt(server); err != nil {
			logger.Fatalf("Failed to apply option: %v", err)
		}
	}
	return server
}

func SetThreadNumber(threadNumber int) func(*app) error {
	return func(a *app) error {
		if threadNumber <= 0 {
			return fmt.Errorf("threadNumber must be greater than 0")
		}
		a.internalThreadNumber = threadNumber
		return nil
	}
}
func SetDataStore(dataStore *datastore.DataStore) func(*app) error {
	return func(a *app) error {
		a.DataStore = dataStore
		return nil
	}
}
func NewEtcdClient(svcEndpoint []string) func(*app) error {
	return func(a *app) error {
		cfg := etcd.Config{
			Endpoints:   svcEndpoint,
			DialTimeout: 5 * time.Second,
		}
		client, err := etcd.New(cfg)
		a.cli = client
		return err
	}
}
func WithContext(ctx context.Context) func(*app) error {
	return func(a *app) error {
		a.Ctx, a.Cancel = context.WithCancel(ctx)
		return nil
	}
}
func UseSignalHandler(handler func()) func(*app) error {
	return func(a *app) error {
		go handler()
		return nil
	}
}
func UseTls(tls *crypto.Tls) func(*app) error {
	return func(a *app) error {
		a.Tls = tls
		return nil
	}
}
func GrpcServer(port string) func(*app) error {
	return func(a *app) error {
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
				return err
			}
			a.Server = grpc.NewServer(
				grpcTlsOption,
				grpc.KeepaliveParams(kasp),
				grpc.KeepaliveEnforcementPolicy(kaep),
				grpc.MaxRecvMsgSize(1024*1024*100), // 100 MB
				grpc.MaxSendMsgSize(1024*1024*100), // 100 MB
				grpc.ChainUnaryInterceptor(
					//middlewares.UnaryServerInterceptor,
					middleware.GetLoggerInterceptor(),
					middleware.GetRecovery(middleware.GetRegistry(middleware.GetSrvMetrics())),
				),
			)
		} else {
			a.Server = grpc.NewServer(
				grpc.KeepaliveParams(kasp),
				grpc.KeepaliveEnforcementPolicy(kaep),
				grpc.MaxRecvMsgSize(1024*1024*100), // 100 MB
				grpc.MaxSendMsgSize(1024*1024*100), // 100 MB
				grpc.ChainUnaryInterceptor(
					//middlewares.UnaryServerInterceptor,
					middleware.GetLoggerInterceptor(),
					middleware.GetRecovery(middleware.GetRegistry(middleware.GetSrvMetrics())),
				),
			)
		}
		svcpb.RegisterDataStoreServiceServer(a.Server, a)
		svcpb.RegisterDataStoreDebugServiceServer(a.Server, a.debug)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			return err
		}
		go func() {
			if err := a.Server.Serve(lis); err != nil {
				logger.Fatalf("Failed to serve: %v", err)
			}
		}()
		return nil
	}
}
