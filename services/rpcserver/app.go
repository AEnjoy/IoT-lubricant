package rpcserver

import (
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
)

type app struct {
	*datastore.DataStore
	*crypto.Tls
	port string

	*grpcServer
}

func (a *app) Run() error {
	a.grpcServer = &grpcServer{PbCoreServiceImpl: PbCoreServiceImpl{
		DataStore: a.DataStore,
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
