package log

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/gin-gonic/gin"
)

type Api struct {
	*datastore.DataStore
}
