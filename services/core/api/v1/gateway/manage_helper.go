package gateway

import (
	"net/http"

	"github.com/AEnjoy/IoT-lubricant/internal/pkg/ssh"
	"github.com/AEnjoy/IoT-lubricant/pkg/form/request"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/crypto"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	helper2 "github.com/AEnjoy/IoT-lubricant/services/core/api/v1/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a Api) getGatewayHostModel(c *gin.Context) (*crypto.Tls, *model.GatewayHost) {
	var req request.AddGatewayRequest
	err := c.BindJSON(&req)
	if err != nil {
		helper2.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.ErrorBind), c)
		return nil, nil
	}
	userInfo, err := helper2.GetClaims(c)
	if err != nil {
		helper2.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.ErrorGetClaimsFailed), c)
		return nil, nil
	}
	if req.PrivateKey == "" && req.PassWd == "" {
		if req.CustomPrivateKey {
			helper2.FailedWithJson(http.StatusBadRequest,
				exception.ErrNewException(err, exceptCode.ErrorGatewayHostNeedPasswdOrPrivateKey), c)
			return nil, nil
		} else {
			key, err := ssh.GetLocalSSHPublicKey()
			if err != nil || key == "" {
				helper2.FailedWithJson(http.StatusTooEarly,
					exception.ErrNewException(err, exceptCode.ErrorGatewayHostNeedPasswdOrPrivateKey,
						exception.WithMsg("failed to get local ssh public key. maybe is not created?"),
						exception.WithMsg("failed to add gateway host due to invalid auth method"),
					), c)
				return nil, nil
			}
			req.PrivateKey = key
		}
	}

	return req.TlsConfig, &model.GatewayHost{
		UserID: userInfo.ID,
		HostID: uuid.NewString(),

		Description: req.Description,
		Host:        req.Host,
		UserName:    req.UserName,
		PassWd:      req.PassWd,
		PrivateKey:  req.PrivateKey,
	}
}

func (a Api) getGatewayRemoveModel(c *gin.Context) *request.RemoveGatewayRequest {
	return helper2.RequestBind[request.RemoveGatewayRequest](c)
}
