package data

import (
	"sync"

	"github.com/aenjoy/iot-lubricant/pkg/utils/crontab"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
)

var dataSendQueue chan *corepb.Data

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
		dataSendQueue = make(chan *corepb.Data, 100)
		_ = crontab.RegisterCron(handelData, "@every 2m")
	})
}
func GetDataSendQueue() <-chan *corepb.Data {
	return dataSendQueue
}
