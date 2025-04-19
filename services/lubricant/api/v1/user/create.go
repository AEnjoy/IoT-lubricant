package user

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/gin-gonic/gin"
)

func (u Api) Create(c *gin.Context) {
	panic("implement me")
}

func (u Api) GetUserInfo(c *gin.Context) {
	claims, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedByServer(err, c)
		return
	}
	helper.SuccessJson(claims.User, c)
}
