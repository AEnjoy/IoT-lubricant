package mq

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
)

type MqService struct {
	mq.Mq[[]byte]
	*datastore.DataStore
}
