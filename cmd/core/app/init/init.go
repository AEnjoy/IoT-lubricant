package init

import (
	"github.com/AEnjoy/IoT-lubricant/cmd/core/app/config"
	_ "github.com/AEnjoy/IoT-lubricant/cmd/core/app/config" // App config
	data "github.com/AEnjoy/IoT-lubricant/cmd/core/app/datastore"
	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
)

// AppInit 初始化App:
//
//	不使用 func init()的原因:避免非关键依赖加载
func AppInit() error {
	c := ioc.Controller.Get(config.APP_NAME).(*config.Config)
	if c.RedisEnable {
		ioc.Controller.Registry(data.APP_NAME, &data.DataStore{CacheEnable: true})
	} else {
		ioc.Controller.Registry(data.APP_NAME, &data.DataStore{})
	}
	return ioc.Controller.Init()
}
