package router

import (
	groups2 "github.com/AEnjoy/IoT-lubricant/services/core/router/groups"
	"github.com/gin-gonic/gin"
)

var CommonGroups = []CommonRouter{
	groups2.UserRoute{},
	groups2.GatewayRoute{},
}

type CommonRouter interface {
	InitRouter(router *gin.RouterGroup)
}
