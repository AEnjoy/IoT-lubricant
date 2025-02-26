package internal

import (
	"sync"

	"github.com/aenjoy/iot-lubricant/pkg/cache"
	"github.com/aenjoy/iot-lubricant/services/lubricant"
	"github.com/aenjoy/iot-lubricant/services/lubricant/auth"
	"github.com/aenjoy/iot-lubricant/services/lubricant/config"
	data "github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	ioc "github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
	mqService "github.com/aenjoy/iot-lubricant/services/lubricant/mq"
	"github.com/aenjoy/iot-lubricant/services/lubricant/repo"
	"github.com/aenjoy/iot-lubricant/services/lubricant/router"
	service2 "github.com/aenjoy/iot-lubricant/services/lubricant/services"
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
			ioc.APP_NAME_CORE_GRPC_AUTH_INTERCEPTOR: &auth.InterceptorImpl{},
			ioc.APP_NAME_CORE_GRPC_SERVER:           &lubricant.Grpc{},
			ioc.APP_NAME_CORE_GATEWAY_SERVICE:       &service2.GatewayService{},
			ioc.APP_NAME_CORE_GATEWAY_AGENT_SERVICE: &service2.AgentService{},
			ioc.APP_NAME_CORE_WEB_SERVER:            &router.WebService{},
			ioc.APP_NAME_CORE_Internal_MQ_SERVICE:   &mqService.MqService{},
		}

		ioc.Controller.LoadObject(objects)
	})

	return ioc.Controller.Init()
}
