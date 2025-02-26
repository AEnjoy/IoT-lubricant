package gateway

import (
	"github.com/AEnjoy/IoT-lubricant/services/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/services/core/services"
)

type Api struct {
	*datastore.DataStore
	services.IGatewayService
	services.IAgentService
}
