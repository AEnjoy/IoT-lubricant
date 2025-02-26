package groups

import (
	"github.com/AEnjoy/IoT-lubricant/services/core/api/v1"
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
}

func (UserRoute) InitRouter(router *gin.RouterGroup) {
	user := router.Group("/user")
	controller := v1.NewUser()

	user.POST("/create", controller.Create)
	user.GET("/info", controller.GetUserInfo)
}
