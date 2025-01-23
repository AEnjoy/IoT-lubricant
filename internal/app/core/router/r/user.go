package r

import (
	v1 "github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1"
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
}

func (UserRoute) InitRouter(router *gin.RouterGroup) {
	user := router.Group("/user")
	controller := v1.NewUser()

	user.POST("/create", controller.Create)
}
