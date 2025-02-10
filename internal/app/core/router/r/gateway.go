package r

import (
	v1 "github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1"
	"github.com/gin-gonic/gin"
)

type GatewayRoute struct {
}

func (GatewayRoute) InitRouter(router *gin.RouterGroup, mids ...gin.HandlerFunc) {
	gateway := router.Group("/gateway", mids...)
	controller := v1.NewGateway()

	gateway.POST("/add-host", controller.AddHost)
	gateway.POST("/host", controller.AddHost)

	gateway.POST("/internal/add-gateway", controller.AddGatewayInternal)
	gateway.POST("/internal/gateway", controller.AddGatewayInternal)
	gateway.POST("/internal/remove-gateway", controller.RemoveGatewayInternal)
	gateway.DELETE("/internal/gateway", controller.RemoveGatewayInternal)

	gateway.POST("/:gatewayid/agent/internal/push-task", controller.AgentPushTask)
	gateway.POST("/:gatewayid/agent/internal/task", controller.AgentPushTask)
}
