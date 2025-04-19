package gateway

import (
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/model/request"
	"github.com/aenjoy/iot-lubricant/pkg/ssh"
	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a Api) getGatewayHostModel(c *gin.Context) (*crypto.Tls, *model.GatewayHost) {
	req := helper.RequestBind[request.AddGatewayRequest](c)
	if req == nil {
		return nil, nil
	}

	userInfo, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.ErrorGetClaimsFailed), c)
		return nil, nil
	}
	if req.PrivateKey == "" && req.PassWd == "" {
		if req.CustomPrivateKey {
			helper.FailedWithJson(http.StatusBadRequest,
				exception.ErrNewException(err, exceptionCode.ErrorGatewayHostNeedPasswdOrPrivateKey), c)
			return nil, nil
		} else {
			key, err := ssh.GetLocalSSHPrivateKey()
			if err != nil || key == "" {
				helper.FailedWithJson(http.StatusTooEarly,
					exception.ErrNewException(err, exceptionCode.ErrorGatewayHostNeedPasswdOrPrivateKey,
						exception.WithMsg("failed to get local ssh public key. maybe is not created?"),
						exception.WithMsg("failed to add gateway host due to invalid auth method"),
					), c)
				return nil, nil
			}
			req.PrivateKey = key
		}
	}

	return req.TlsConfig, &model.GatewayHost{
		UserID: userInfo.User.Id,
		HostID: uuid.NewString(),

		Description: req.Description,
		Host:        req.Host,
		UserName:    req.UserName,
		PassWd:      req.PassWd,
		PrivateKey:  req.PrivateKey,
	}
}

func (a Api) getGatewayRemoveModel(c *gin.Context) *request.RemoveGatewayRequest {
	return helper.RequestBind[request.RemoveGatewayRequest](c)
}
