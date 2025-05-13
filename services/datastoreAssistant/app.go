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
	subscribed  map[string]struct{} //projectId-subscribed
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

	a.subscribed = make(map[string]struct{})
	regPool, err := ants.NewPoolWithFunc(a.internalThreadNumber, a.handel, ants.WithPreAlloc(true))
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
