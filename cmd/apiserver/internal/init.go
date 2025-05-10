package internal

import (
	"sync"

	"github.com/aenjoy/iot-lubricant/services/apiserver/router"
	"github.com/aenjoy/iot-lubricant/services/apiserver/services"
	backendService "github.com/aenjoy/iot-lubricant/services/apiserver/services/backend"
	"github.com/aenjoy/iot-lubricant/services/corepkg/cache"
	"github.com/aenjoy/iot-lubricant/services/corepkg/config"
	data "github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
	"github.com/aenjoy/iot-lubricant/services/corepkg/logCollect"
	mqService "github.com/aenjoy/iot-lubricant/services/corepkg/mq"
	"github.com/aenjoy/iot-lubricant/services/corepkg/repo"
	"github.com/aenjoy/iot-lubricant/services/corepkg/syncQueue"
)

var once sync.Once

// AppInit 初始化App:
//
//	不使用 func init()的原因:避免非关键依赖加载
func AppInit() error {
	// AppObjects witch will be registered with default option
	once.Do(func() {
		var objects = map[string]ioc.Object{
			config.APP_NAME: config.GetConfig(),

			ioc.APP_NAME_CORE_CACHE:                 &cache.RedisCli[string]{},
			ioc.APP_NAME_CORE_DATABASE:              &repo.CoreDb{},
			ioc.APP_NAME_CORE_DATABASE_STORE:        &data.DataStore{CacheEnable: config.GetConfig().RedisEnable},
			ioc.APP_NAME_CORE_GATEWAY_SERVICE:       &services.GatewayService{},
			ioc.APP_NAME_CORE_GATEWAY_AGENT_SERVICE: &services.AgentService{},
			ioc.APP_NAME_CORE_WEB_SERVER:            &router.WebService{},
			ioc.APP_NAME_CORE_PROJECT_SERVICE:       &services.ProjectService{},

			ioc.APP_NAME_CORE_Internal_MQ_SERVICE:         &mqService.MqService{},
			ioc.APP_NAME_CORE_Internal_SyncTask_SERVICE:   &syncQueue.SyncTaskQueue{},
			ioc.APP_NAME_CORE_Internal_LOGGER_SERVICE:     &logCollect.Log{},
			ioc.APP_NAME_CORE_Internal_Handler_DataUpload: &backendService.DataHandler{},
			ioc.APP_NAME_CORE_Internal_Handler_ErrLogs:    &backendService.ErrLogCollect{},
		}

		ioc.Controller.LoadObject(objects)
	})

	return ioc.Controller.Init()
}
