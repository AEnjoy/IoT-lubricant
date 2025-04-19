package log

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
)

type Api struct {
	*datastore.DataStore
}
