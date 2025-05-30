package router

import (
	"github.com/aenjoy/iot-lubricant/services/apiserver/router/groups"
	"github.com/gin-gonic/gin"
)

var CommonGroups = []CommonRouter{
	groups.UserRoute{},
	groups.GatewayRoute{},
	groups.MonitorRoute{},
	groups.AgentRoute{},
	groups.TaskRoute{},
	groups.ProjectRoute{},
}

type CommonRouter interface {
	InitRouter(router *gin.RouterGroup)
}
