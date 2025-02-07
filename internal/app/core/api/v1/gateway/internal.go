package gateway

import (
	"net/http"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1/helper"
	"github.com/AEnjoy/IoT-lubricant/internal/model/form/response"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a Api) AddGatewayInternal(c *gin.Context) {
	tlsConfig, gatewayHostInfo := a.getGatewayHostModel(c)
	if gatewayHostInfo == nil {
		return
	}
	err := a.IGatewayService.AddHostInternal(c, gatewayHostInfo)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.AddGatewayHostFailed), c)
		return
	}
	gatewayid := uuid.NewString()
	err = a.IGatewayService.AddGatewayInternal(c, gatewayHostInfo.UserID, gatewayid, gatewayHostInfo.Description, tlsConfig)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.AddGatewayFailed), c)
		return
	}
	helper.SuccessJson(response.AddGatewayResponse{GatewayID: gatewayid}, c)
}
func (a Api) RemoveGatewayInternal(c *gin.Context) {
	req := a.getGatewayRemoveModel(c)
	if req == nil {
		return
	}

	if err := a.IGatewayService.RemoveGatewayInternal(c, req.GatewayID); err != nil {
		helper.FailedWithJson(http.StatusInternalServerError, err.(*exception.Exception), c)
		return
	}
	helper.SuccessJson(nil, c)
}
