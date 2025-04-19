package gateway

import (
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/model/request"
	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/gin-gonic/gin"
)

func (a Api) DescriptionGateway(c *gin.Context) {
	claims, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.New(exceptionCode.ErrorGetClaimsFailed, exception.WithMsg("claims is empty")), c)
		return
	}
	userid := claims.User.Id
	gatewayid := c.Param("gatewayid")
	gatewayinfo, err := a.IGatewayService.DescriptionGateway(c, userid, gatewayid)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.GetGatewayFailed), c)
		return
	}
	helper.SuccessJson(gatewayinfo, c)
}

func (a Api) EditGateway(c *gin.Context) {
	req := helper.RequestBind[request.EditGatewayRequest](c)
	if req == nil {
		return
	}

	err := a.IGatewayService.EditGateway(c, req.GatewayID, req.Description, req.TlsConfig)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.DbUpdateGatewayInfoFailed), c)
		return
	}
	helper.SuccessJson(response.Empty{}, c)
}
