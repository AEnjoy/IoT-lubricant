package groups

import (
	v1 "github.com/aenjoy/iot-lubricant/services/lubricant/api"
	"github.com/gin-gonic/gin"
)

type AgentRoute struct{}

func (AgentRoute) InitRouter(router *gin.RouterGroup) {
	agent := router.Group("/agent")
	controller := v1.NewAgent()

	agent.GET("/operator", controller.Operator)
	agent.POST("/set", controller.Set)
	agent.GET("/get-data", controller.GetData)
	agent.GET("/info", controller.GetAgentInfo)
	agent.GET("/list", controller.List)
}
