package router

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1/helper"
	"github.com/gin-gonic/gin"
)

func Health(ctx *gin.Context) {
	helper.SuccessJson("ok", ctx)
}
