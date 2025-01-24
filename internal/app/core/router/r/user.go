package r

import (
	v1 "github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1"
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
}

func (UserRoute) InitRouter(router *gin.RouterGroup, mids ...gin.HandlerFunc) {
	user := router.Group("/user", mids...)
	controller := v1.NewUser()

	user.POST("/create", controller.Create)
}
