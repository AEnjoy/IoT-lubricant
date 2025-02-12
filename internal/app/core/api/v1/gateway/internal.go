package gateway

import (
	"encoding/base64"
	"net/http"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1/helper"
	"github.com/AEnjoy/IoT-lubricant/internal/model/form/request"
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
				exception.ErrNewException(err, exceptCode.ErrorDecodeFailed,
					exception.WithMsg("error in openapi doc file")), c)
			return
		}
	}
	if req.EnableConf != "" {
		enableFile, err = base64.StdEncoding.DecodeString(req.EnableConf)
		if err != nil {
			helper.FailedWithJson(http.StatusInternalServerError,
				exception.ErrNewException(err, exceptCode.ErrorDecodeFailed,
					exception.WithMsg("error in openapi enable config file")), c)
			return
		}
	}

	var taskID *string
	agentID, err := a.IGatewayService.AddAgentInternal(c, taskID, gatewayID, req, openapidoc, enableFile)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError, err.(*exception.Exception), c)
		return
	}
	helper.SuccessJson(response.AddAgentResponse{AgentID: agentID,
		PushAgentTaskResponse: response.PushAgentTaskResponse{TaskID: *taskID},
	}, c)
}
