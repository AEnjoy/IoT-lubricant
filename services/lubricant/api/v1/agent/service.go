package agent

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/services"
)

type Api struct {
	*datastore.DataStore
	services.IAgentService
}
