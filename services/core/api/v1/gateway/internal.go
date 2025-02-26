package gateway

import (
	"encoding/base64"
	"net/http"

	"github.com/AEnjoy/IoT-lubricant/pkg/form/request"
	response2 "github.com/AEnjoy/IoT-lubricant/pkg/form/response"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	helper2 "github.com/AEnjoy/IoT-lubricant/services/core/api/v1/helper"
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
		helper2.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.AddGatewayHostFailed), c)
		return
	}
	gatewayid := uuid.NewString()
	err = a.IGatewayService.AddGatewayInternal(c, gatewayHostInfo.UserID, gatewayid, gatewayHostInfo.Description, tlsConfig)
	if err != nil {
		helper2.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.AddGatewayFailed), c)
		return
	}
	helper2.SuccessJson(response2.AddGatewayResponse{GatewayID: gatewayid}, c)
}
func (a Api) RemoveGatewayInternal(c *gin.Context) {
	req := a.getGatewayRemoveModel(c)
	if req == nil {
		return
	}

	if err := a.IGatewayService.RemoveGatewayInternal(c, req.GatewayID); err != nil {
		helper2.FailedWithJson(http.StatusInternalServerError, err.(*exception.Exception), c)
		return
	}
	helper2.SuccessJson(nil, c)
}

func (a Api) AddAgentInternal(c *gin.Context) {
	req := helper2.RequestBind[request.AddAgentRequest](c)
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
			helper2.FailedWithJson(http.StatusInternalServerError,
				exception.ErrNewException(err, exceptCode.ErrorDecodeFailed,
					exception.WithMsg("error in openapi doc file")), c)
			return
		}
	}
	if req.EnableConf != "" {
		enableFile, err = base64.StdEncoding.DecodeString(req.EnableConf)
		if err != nil {
			helper2.FailedWithJson(http.StatusInternalServerError,
				exception.ErrNewException(err, exceptCode.ErrorDecodeFailed,
					exception.WithMsg("error in openapi enable config file")), c)
			return
		}
	}

	_id := uuid.NewString()
	var taskID = &_id

	agentID, err := a.IGatewayService.AddAgentInternal(c, taskID, gatewayID, req, openapidoc, enableFile)
	if err != nil {
		helper2.FailedWithJson(http.StatusInternalServerError, err.(*exception.Exception), c)
		return
	}
	helper2.SuccessJson(response2.AddAgentResponse{AgentID: agentID,
		PushAgentTaskResponse: response2.PushAgentTaskResponse{TaskID: *taskID},
	}, c)
}
