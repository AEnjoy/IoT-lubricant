package datastoreAssistant

import (
	"strconv"
	"sync"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/panjf2000/ants/v2"
	etcd "go.etcd.io/etcd/client/v3"
)

type app struct {
	subscribed map[string]struct{} //userId-subscribed
	l          sync.Mutex
	*datastore.DataStore
	cli                  *etcd.Client
	internalThreadNumber int
}

func (a *app) Run() error {
	regPool, err := ants.NewPoolWithFunc(a.internalThreadNumber, a.register, ants.WithPreAlloc(true))
	if err != nil {
		return err
	}
	a.Mq.SubscribeBytes()
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
func SetThreadNumber(threadNumber string) func(*app) error {
	return func(a *app) error {
		var err error
		a.internalThreadNumber, err = strconv.Atoi(threadNumber)
		return err
	}
}
func SetDataStore(dataStore *datastore.DataStore) func(*app) error {
	return func(a *app) error {
		a.DataStore = dataStore
		return nil
	}
}
func NewEtcdClient(svcEndpoint string) func(*app) error {
	return func(a *app) error {
		cfg := etcd.Config{
			Endpoints:   []string{svcEndpoint},
			DialTimeout: 5 * time.Second,
		}
		client, err := etcd.New(cfg)
		a.cli = client
		return err
	}
}
