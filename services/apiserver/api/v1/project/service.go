package project

import (
	"github.com/aenjoy/iot-lubricant/services/apiserver/services"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
)

type Api struct {
	*datastore.DataStore
	services.IProjectService
}
