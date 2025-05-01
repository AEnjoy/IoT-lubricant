package rpcserver

import (
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/panjf2000/ants/v2"
)

type app struct {
	*datastore.DataStore
	*crypto.Tls
	port string

	*grpcServer
}

func (a *app) Run() error {
	pool, err := ants.NewPool(2048, ants.WithPreAlloc(true))
	if err != nil {
		return err
	}
	a.grpcServer = &grpcServer{PbCoreServiceImpl: PbCoreServiceImpl{
		DataStore: a.DataStore,
		pool:      pool,
	}}
	serverStart, err := a.grpcInit()
	if err != nil {
		return err
	}

	return serverStart()
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
func SetPort(port string) func(*app) error {
	return func(s *app) error {
		s.port = port
		return nil
	}
}

func UseTls(tls *crypto.Tls) func(*app) error {
	return func(a *app) error {
		a.Tls = tls
		return nil
	}
}
func UseDataStore(store *datastore.DataStore) func(*app) error {
	return func(a *app) error {
		a.DataStore = store
		return nil
	}
}
func UseSignalHandler(handler func()) func(*app) error {
	return func(a *app) error {
		go handler()
		return nil
	}
}
