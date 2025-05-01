package datastoreAssistant

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/datastoreAssistant/api"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"github.com/panjf2000/ants/v2"

	etcd "go.etcd.io/etcd/client/v3"
)

type app struct {
	subscribed  map[string]struct{} //userId-subscribed
	subMapMutex sync.Mutex

	*datastore.DataStore
	cli                  *etcd.Client
	internalThreadNumber int
	pool                 *ants.PoolWithFunc

	StoringMutex sync.Mutex // 只要有任意的线程在存储数据，则不允许未消费完就退出

	Ctx    context.Context
	Cancel context.CancelFunc
}

func (a *app) Run() error {
	_app = a
	regPool, err := ants.NewPoolWithFunc(a.internalThreadNumber, a.handel, ants.WithPreAlloc(true))
	if err != nil {
		return err
	}
	a.pool = regPool

	hostname, _ := os.Hostname()
	consumerID := fmt.Sprintf("%s-%d", hostname, time.Now().UnixNano())
	userIdBytesCh, err := a.V2mq.Subscribe(constant.DATASTORE_USER)
	if err != nil {
		return err
	}
	_, err = a.registerConsumer(consumerID)
	if err != nil {
		return fmt.Errorf("[%s] Failed to register consumer: %v", consumerID, err)
	}
	go a.watchAssignments(consumerID)
	go a.campaignForLeadership(consumerID)

	for {
		select {
		case <-a.Ctx.Done():
			time.Sleep(3 * time.Second)
			os.Exit(0)
		case userIdBytes := <-userIdBytesCh:
			userID := string(userIdBytes.([]byte))
			partitionID := api.GetPartition(userID, 64)
			if isPartitionAssigned(partitionID) {
				a.subMapMutex.Lock()
				_, ok := a.subscribed[userID]
				if !ok {
					err := a.pool.Invoke(userID)
					if err != nil {
						logg.L.Errorf("[%s] Failed to invoke pool: %v", userID, err)
					} else {
						a.subscribed[userID] = struct{}{}
					}
				}
				a.subMapMutex.Unlock()
			}
		}
	}
}

var _app *app

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
