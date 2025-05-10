package internal

import (
	"sync"

	"github.com/aenjoy/iot-lubricant/services/corepkg/config"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
	"github.com/aenjoy/iot-lubricant/services/corepkg/logCollect"
	"github.com/aenjoy/iot-lubricant/services/corepkg/mq"
	"github.com/aenjoy/iot-lubricant/services/corepkg/repo"
)

var once sync.Once

// AppInit 初始化App:
func AppInit() error {
	// AppObjects which will be registered with default option
	once.Do(func() {
		var objects = map[string]ioc.Object{
			config.APP_NAME: config.GetConfig(),

			ioc.APP_NAME_CORE_DATABASE:       &repo.CoreDb{},
			ioc.APP_NAME_CORE_DATABASE_STORE: &datastore.DataStore{},

			ioc.APP_NAME_CORE_Internal_MQ_SERVICE:     &mq.MqService{},
			ioc.APP_NAME_CORE_Internal_LOGGER_SERVICE: &logCollect.Log{},
		}

		ioc.Controller.LoadObject(objects)
	})

	return ioc.Controller.Init()
}
