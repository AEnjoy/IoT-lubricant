package router

import (
	"github.com/gin-gonic/gin"
)

func CommonGroups() []CommonRouter {
	return []CommonRouter{}
}

type CommonRouter interface {
	InitRouter(router *gin.RouterGroup)
}
