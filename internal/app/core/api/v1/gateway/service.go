package gateway

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/service"
)

type Api struct {
	*datastore.DataStore
	service.IGatewayService
}
