package groups

import (
	v1 "github.com/aenjoy/iot-lubricant/services/lubricant/api"
	"github.com/gin-gonic/gin"
)

type MonitorRoute struct {
}

func (MonitorRoute) InitRouter(router *gin.RouterGroup) {
	m := router.Group("/monitor")
	controller := v1.NewMonitor()

	m.GET("/info", controller.BaseInfo)
}
