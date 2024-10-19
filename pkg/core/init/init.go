package init

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/auth"
	"github.com/AEnjoy/IoT-lubricant/pkg/core"
	"github.com/AEnjoy/IoT-lubricant/pkg/core/config"
	data "github.com/AEnjoy/IoT-lubricant/pkg/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
)

// AppInit 初始化App:
//
//	不使用 func init()的原因:避免非关键依赖加载
func AppInit() error {
	// AppObjects witch will be registered with default option
	var objects = map[string]ioc.Object{
		config.APP_NAME:            ioc.Controller.Get(config.APP_NAME).(*config.Config),
		ioc.APP_NAME_CORE_DATABASE: &model.CoreDb{},
		ioc.APP_NAME_CORE_DATABASE_STORE: func() ioc.Object {
			c := ioc.Controller.Get(config.APP_NAME).(*config.Config)
			if c.RedisEnable {
				return &data.DataStore{CacheEnable: true}
			}
			return &data.DataStore{}
		}(),
		ioc.APP_NAME_CORE_GRPC_AUTH_INTERCEPTOR: &auth.InterceptorImpl{},
		ioc.APP_NAME_CORE_GRPC_SERVER:           &core.Grpc{},
	}

	ioc.Controller.LoadObject(objects)
	return ioc.Controller.Init()
}
