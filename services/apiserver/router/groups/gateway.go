package groups

import (
	v1 "github.com/aenjoy/iot-lubricant/services/apiserver/api"
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

	gateway.POST("/add-gateway", controller.AddGateway)
	gateway.POST("/add", controller.AddGateway)
	gateway.POST("/remove-gateway", controller.RemoveGateway)
	gateway.DELETE("/gateway", controller.RemoveGateway)

	gateway.POST("/:gatewayid/agent/internal/push-task", controller.AgentPushTask)
	gateway.POST("/:gatewayid/agent/internal/task", controller.AgentPushTask)
	gateway.POST("/:gatewayid/agent/internal/add", controller.AddAgentInternal)
	gateway.GET("/:gatewayid/description", controller.DescriptionGateway)
	gateway.POST("/:gatewayid/edit", controller.EditGateway)
}
