package gateway

import (
	"encoding/base64"
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/model/request"
	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/xid"
)

func (a Api) AddGatewayInternal(c *gin.Context) {
	tlsConfig, gatewayHostInfo := a.getGatewayHostModel(c)
	if gatewayHostInfo == nil {
		return
	}
	err := a.IGatewayService.AddHostInternal(c, gatewayHostInfo)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.AddGatewayHostFailed), c)
		return
	}

	gatewayid, ok := c.GetQuery("gateway-id")
	if !ok {
		gatewayid = xid.New().String()
	}

	err = a.IGatewayService.AddGatewayInternal(c, gatewayHostInfo.UserID, gatewayid, gatewayHostInfo.Description, tlsConfig)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.AddGatewayFailed), c)
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

func (a Api) AddAgentInternal(c *gin.Context) {
	req := helper.RequestBind[request.AddAgentRequest](c)
	if req == nil {
		return
	}
	gatewayID := c.Param("gatewayid")

	var (
		openapidoc []byte
		enableFile []byte

		err error
	)
	if req.OpenApiDoc != "" {
		openapidoc, err = base64.StdEncoding.DecodeString(req.OpenApiDoc)
		if err != nil {
			helper.FailedWithJson(http.StatusInternalServerError,
				exception.ErrNewException(err, exceptionCode.ErrorDecodeFailed,
					exception.WithMsg("error in openapi doc file")), c)
			return
		}
	}
	if req.EnableConf != "" {
		enableFile, err = base64.StdEncoding.DecodeString(req.EnableConf)
		if err != nil {
			helper.FailedWithJson(http.StatusInternalServerError,
				exception.ErrNewException(err, exceptionCode.ErrorDecodeFailed,
					exception.WithMsg("error in openapi enable config file")), c)
			return
		}
	}
	claims, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.New(exceptionCode.ErrorGetClaimsFailed, exception.WithMsg("claims is empty")), c)
		return
	}
	userid := claims.User.Id
	_id := uuid.NewString()
	var taskID = &_id

	agentID, err := a.IGatewayService.AddAgentInternal(c, taskID, userid, gatewayID, req, openapidoc, enableFile)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError, err.(*exception.Exception), c)
		return
	}
	helper.SuccessJson(response.AddAgentResponse{AgentID: agentID,
		PushAgentTaskResponse: response.PushAgentTaskResponse{TaskID: *taskID},
	}, c)
}
