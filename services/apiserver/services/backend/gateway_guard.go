package backend

import (
	"context"
	"fmt"
	"strings"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
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
		go func(str string) {
			// str is `"%s<!SPLIT!>%s", userid, gatewayid`
			var userid, gatewayid string
			result := strings.Split(str, "<!SPLIT!>")
			if len(result) == 2 {
				userid = result[0]
				gatewayid = result[1]
			} else {
				logger.Errorf("internalError: failed to split gateway id: %s", str)
				return
			}

			txn := g.dataCli.ICoreDb.Begin()
			err := g.dataCli.ICoreDb.SetGatewayStatus(context.Background(), txn, userid, gatewayid, "offline")
			g.dataCli.ICoreDb.Commit(txn)
			if err != nil {
				logger.Errorf("failed to set gateway status: %v", err)
			}
		}(string(ch.([]byte)))
	}
}
