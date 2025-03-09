package lubricant

import (
	"context"
	"fmt"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
)

func gatewayStatusGuard() {
	time.Sleep(3 * time.Second)
	dataStore := dataCli()
	gatewayCh, err := dataStore.Mq.SubscribeBytes(fmt.Sprintf("/monitor/%s/offline", taskTypes.TargetGateway))
	if err != nil {
		return
	}
	for ch := range gatewayCh {
		go func(id string) {
			txn := dataStore.ICoreDb.Begin()
			err := dataStore.ICoreDb.SetGatewayStatus(context.Background(), txn, id, "offline")
			dataStore.ICoreDb.Commit(txn)
			if err != nil {
				logger.Errorf("failed to set gateway status: %v", err)
			}
		}(string(ch.([]byte)))
	}
}
