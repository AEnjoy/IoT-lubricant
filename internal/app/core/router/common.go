package router

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/router/r"
	"github.com/gin-gonic/gin"
)

func CommonGroups() []CommonRouter {
	return []CommonRouter{
		r.UserRoute{},
		r.GatewayRoute{},
	}
}

type CommonRouter interface {
	InitRouter(router *gin.RouterGroup, mids ...gin.HandlerFunc)
}
