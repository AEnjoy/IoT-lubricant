package router

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/gin-gonic/gin"
)

func Health(ctx *gin.Context) {
	helper.SuccessJson("ok", ctx)
}
