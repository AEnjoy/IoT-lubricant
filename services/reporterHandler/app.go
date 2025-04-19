package reporterHandler

import (
	"fmt"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/syncQueue"
	"github.com/panjf2000/ants/v2"
)

type app struct {
	*datastore.DataStore
	*syncQueue.SyncTaskQueue
}

func (a app) Run() error {
	reporterPool, err := ants.NewPoolWithFunc(512, a.reporterPayload, ants.WithPreAlloc(true))
	if err != nil {
		return err
	}
	guardPool, err := ants.NewPoolWithFunc(64, a.gatewayOfflineGuardPayload, ants.WithPreAlloc(true))
	if err != nil {
		return err
	}
	reportDataCh, err := a.DataStore.SubscribeBytes("/handler/report")
	if err != nil {
		return err
	}
	gatewayCh, err := a.DataStore.Mq.SubscribeBytes(fmt.Sprintf("/monitor/%s/offline", taskTypes.TargetGateway))
	if err != nil {
		return err
	}
	for {
		select {
		case data := <-reportDataCh:
			_ = reporterPool.Invoke(data)
		case data := <-gatewayCh:
			_ = guardPool.Invoke(data)
		}
	}
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
func UseSignalHandler(handler func()) func(*app) error {
	return func(a *app) error {
		go handler()
		return nil
	}
}
func UseDataStore(store *datastore.DataStore) func(*app) error {
	return func(a *app) error {
		a.DataStore = store
		return nil
	}
}
func UseTaskQueue(queue *syncQueue.SyncTaskQueue) func(*app) error {
	return func(a *app) error {
		a.SyncTaskQueue = queue
		return nil
	}
}
