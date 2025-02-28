package gateway

import (
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/gin-gonic/gin"
)

func (a Api) DescriptionGateway(c *gin.Context) {
	gatewayid := c.Param("gatewayid")
	gatewayinfo, err := a.IGatewayService.DescriptionGateway(c, gatewayid)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.GetGatewayFailed), c)
		return
	}
	helper.SuccessJson(gatewayinfo, c)
}
