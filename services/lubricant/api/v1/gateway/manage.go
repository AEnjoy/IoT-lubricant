package gateway

import (
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/ssh"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/gin-gonic/gin"
)

func (a Api) AddHost(c *gin.Context) {
	_, gatewayHostInfo := a.getGatewayHostModel(c)
	if gatewayHostInfo == nil {
		return
	}
	_, err := ssh.NewSSHClient(gatewayHostInfo, true)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err,
				exceptionCode.LinkToGatewayFailed,
				exception.WithMsg("test linker failed: failed to link to target host")), c)
		return
	}
	err = a.IGatewayService.AddHost(c, gatewayHostInfo)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.AddGatewayHostFailed), c)
		return
	}
	helper.SuccessJson(gatewayHostInfo, c)
}
