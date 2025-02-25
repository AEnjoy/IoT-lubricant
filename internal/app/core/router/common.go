package router

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/router/groups"
	"github.com/gin-gonic/gin"
)

var CommonGroups = []CommonRouter{
	groups.UserRoute{},
	groups.GatewayRoute{},
}

type CommonRouter interface {
	InitRouter(router *gin.RouterGroup)
}
