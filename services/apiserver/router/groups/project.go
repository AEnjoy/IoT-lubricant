package groups

import (
	v1 "github.com/aenjoy/iot-lubricant/services/apiserver/api"
	"github.com/gin-gonic/gin"
)

type ProjectRoute struct{}

func (ProjectRoute) InitRouter(router *gin.RouterGroup) {
	project := router.Group("/project")
	controller := v1.NewProject()

	project.POST("/add", controller.AddProject)
	project.POST("/remove", controller.RemoveProject)
	project.POST("/bind-agent", controller.AgentBindProject)
	project.POST("/bind-washer", controller.AgentBindWasher)
	project.POST("/add-washer", controller.AddWasher)
	project.POST("/add-engine", controller.AddDataStoreEngine)
	project.GET("/engine-status", controller.GetProjectDataStoreEngineStatus)
	project.POST("/update-engine", controller.UpdateEngineInfo)
}
