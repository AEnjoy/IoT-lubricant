package groups

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1"
	"github.com/gin-gonic/gin"
)

type GatewayRoute struct {
}

func (GatewayRoute) InitRouter(router *gin.RouterGroup) {
	gateway := router.Group("/gateway")
	controller := v1.NewGateway()

	gateway.POST("/add-host", controller.AddHost)
	gateway.POST("/host", controller.AddHost)
	gateway.GET("/host", controller.DescriptionHost)
	gateway.GET("/host/description", controller.DescriptionHost)
	gateway.GET("/list-host", controller.ListHosts)

	gateway.POST("/internal/add-gateway", controller.AddGatewayInternal)
	gateway.POST("/internal/gateway", controller.AddGatewayInternal)
	gateway.POST("/internal/remove-gateway", controller.RemoveGatewayInternal)
	gateway.DELETE("/internal/gateway", controller.RemoveGatewayInternal)

	gateway.POST("/:gatewayid/agent/internal/push-task", controller.AgentPushTask)
	gateway.POST("/:gatewayid/agent/internal/task", controller.AgentPushTask)
	gateway.POST("/:gatewayid/agent/internal/add", controller.AddAgentInternal)
}
