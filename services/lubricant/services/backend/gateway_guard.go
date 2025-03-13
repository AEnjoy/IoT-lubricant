package backend

import (
	"context"
	"fmt"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
)

type GatewayGuard struct {
	dataCli *datastore.DataStore
}

func (g *GatewayGuard) Init() error {
	g.dataCli = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	go g.handler()
	return nil
}
func (g *GatewayGuard) handler() {
	gatewayCh, err := g.dataCli.Mq.SubscribeBytes(fmt.Sprintf("/monitor/%s/offline", taskTypes.TargetGateway))
	if err != nil {
		panic(err)
	}
	for ch := range gatewayCh {
		go func(id string) {
			txn := g.dataCli.ICoreDb.Begin()
			err := g.dataCli.ICoreDb.SetGatewayStatus(context.Background(), txn, id, "offline")
			g.dataCli.ICoreDb.Commit(txn)
			if err != nil {
				logger.Errorf("failed to set gateway status: %v", err)
			}
		}(string(ch.([]byte)))
	}
}
