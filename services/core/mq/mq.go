package mq

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/AEnjoy/IoT-lubricant/services/core/datastore"
)

type MqService struct {
	mq.Mq
	*datastore.DataStore
}
