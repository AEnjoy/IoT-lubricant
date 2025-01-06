package data

import (
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/crontab"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
)

var dataSendQueue chan *core.Data

var once sync.Once

func InitDataSendQueue() {
	handelData := func() {
		agentDataStore.Range(func(key, value any) bool {
			id := key.(string)
			api := value.(Apis)
			data := api.Pop()
			if data != nil {
				data.AgentID = id
				dataSendQueue <- data
			}
			return true
		})
	}
	once.Do(func() {
		dataSendQueue = make(chan *core.Data, 100)
		_ = crontab.RegisterCron(handelData, "@every 2m")
	})
}
func GetDataSendQueue() <-chan *core.Data {
	return dataSendQueue
}
