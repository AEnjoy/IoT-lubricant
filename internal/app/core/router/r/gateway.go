package r

import (
	v1 "github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1"
	"github.com/gin-gonic/gin"
)

type GatewayRoute struct {
}

func (GatewayRoute) InitRouter(router *gin.RouterGroup, mids ...gin.HandlerFunc) {
	user := router.Group("/gateway", mids...)
	controller := v1.NewGateway()

	user.POST("/add-host", controller.AddHost)
}
