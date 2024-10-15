package init

import (
	"github.com/AEnjoy/IoT-lubricant/cmd/core/app"
	"github.com/AEnjoy/IoT-lubricant/cmd/core/app/config" // default app config
	data "github.com/AEnjoy/IoT-lubricant/cmd/core/app/datastore"
	"github.com/AEnjoy/IoT-lubricant/pkg/auth"
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
		ioc.APP_NAME_CORE_GRPC_AYTH_INTERCEPTOR: &auth.InterceptorImpl{},
		ioc.APP_NAME_CORE_GRPC_SERVER:           &app.Grpc{},
	}

	ioc.Controller.LoadObject(objects)
	return ioc.Controller.Init()
}
