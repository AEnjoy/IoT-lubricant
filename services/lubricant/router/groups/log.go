package groups

import (
	v1 "github.com/aenjoy/iot-lubricant/services/lubricant/api"
	"github.com/gin-gonic/gin"
)

type LogRoute struct{}

func (LogRoute) InitRouter(router *gin.RouterGroup) {
	log := router.Group("/log")
	controller := v1.NewLog()

	log.GET("/list", controller.GetLogList)
}
