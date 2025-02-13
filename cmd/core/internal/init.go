package internal

import (
	"sync"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core"
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/config"
	data "github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/service"
	"github.com/AEnjoy/IoT-lubricant/internal/cache"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/internal/pkg/auth"
)

var once sync.Once

// AppInit 初始化App:
//
//	不使用 func init()的原因:避免非关键依赖加载
func AppInit() error {
	// AppObjects witch will be registered with default option
	once.Do(func() {
		var objects = map[string]ioc.Object{
			config.APP_NAME:                         config.GetConfig(),
			ioc.APP_NAME_CORE_CACHE:                 &cache.RedisCli[string]{},
			ioc.APP_NAME_CORE_DATABASE:              &repo.CoreDb{},
			ioc.APP_NAME_CORE_DATABASE_STORE:        &data.DataStore{CacheEnable: config.GetConfig().RedisEnable},
			ioc.APP_NAME_CORE_GRPC_AUTH_INTERCEPTOR: &auth.InterceptorImpl{},
			ioc.APP_NAME_CORE_GRPC_SERVER:           &core.Grpc{},
			ioc.APP_NAME_CORE_GATEWAY_SERVICE:       &service.GatewayService{},
			ioc.APP_NAME_CORE_GATEWAY_AGENT_SERVICE: &service.AgentService{},
		}

		ioc.Controller.LoadObject(objects)
	})

	return ioc.Controller.Init()
}
