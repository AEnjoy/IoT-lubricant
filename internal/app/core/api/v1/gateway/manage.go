package gateway

import (
	"net/http"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1/helper"
	"github.com/AEnjoy/IoT-lubricant/internal/pkg/ssh"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
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
				exceptCode.LinkToGatewayFailed,
				exception.WithMsg("test linker failed: failed to link to target host")), c)
		return
	}
	err = a.IGatewayService.AddHost(c, gatewayHostInfo)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.AddGatewayHostFailed), c)
		return
	}
	helper.SuccessJson(gatewayHostInfo, c)
}
