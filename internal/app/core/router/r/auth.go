package r

import (
	v1 "github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1"
	"github.com/gin-gonic/gin"
)

type AuthRoute struct{}

func (AuthRoute) InitRouter(router *gin.RouterGroup) {
	controller := v1.NewAuth()

	router.POST("/signin", controller.Signin)
}
