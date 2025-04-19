package mq

import (
	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
)

type MqService struct {
	mq.Mq
	*datastore.DataStore
}
