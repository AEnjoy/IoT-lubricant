package groups

import (
	v1 "github.com/aenjoy/iot-lubricant/services/apiserver/api"
	"github.com/gin-gonic/gin"
)

type TaskRoute struct{}

func (TaskRoute) InitRouter(router *gin.RouterGroup) {
	task := router.Group("/task")
	controller := v1.NewTask()

	task.GET("/query", controller.QueryTask)
	task.GET("/list", controller.GetTaskList)
}
