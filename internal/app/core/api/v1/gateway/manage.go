package gateway

import (
	"net/http"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1/helper"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/form/request"
	"github.com/AEnjoy/IoT-lubricant/internal/pkg/ssh"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a Api) AddHost(c *gin.Context) {
	var req request.AddGatewayHostRequest
	err := c.BindJSON(&req)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.ErrorBind), c)
		return
	}
	userInfo, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.ErrorGetClaimsFailed), c)
		return
	}
	if req.PrivateKey == "" && req.PassWd == "" {
		if req.CustomPrivateKey {
			helper.FailedWithJson(http.StatusBadRequest,
				exception.ErrNewException(err, exceptCode.ErrorGatewayHostNeedPasswdOrPrivateKey), c)
			return
		} else {
			key, err := ssh.GetLocalSSHPublicKey()
			if err != nil || key == "" {
				helper.FailedWithJson(http.StatusTooEarly,
					exception.ErrNewException(err, exceptCode.ErrorGatewayHostNeedPasswdOrPrivateKey,
						exception.WithMsg("failed to get local ssh public key. maybe is not created?"),
						exception.WithMsg("failed to add gateway host due to invalid auth method"),
					), c)
				return
			}
			req.PrivateKey = key
		}
	}

	gatewayHostInfo := model.GatewayHost{
		UserID: userInfo.ID,
		HostID: uuid.NewString(),

		Description: req.Description,
		Host:        req.Host,
		UserName:    req.UserName,
		PassWd:      req.PassWd,
		PrivateKey:  req.PrivateKey,
	}
	_, err = ssh.NewSSHClient(&gatewayHostInfo, true)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err,
				exceptCode.LinkToGatewayFailed,
				exception.WithMsg("test linker failed: failed to link to target host")), c)
		return
	}
	err = a.IGatewayService.AddHost(c, &gatewayHostInfo)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.AddGatewayHostFailed), c)
		return
	}
	helper.SuccessJson(gatewayHostInfo, c)
}
