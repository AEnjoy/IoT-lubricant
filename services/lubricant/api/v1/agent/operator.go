package agent

import (
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/gin-gonic/gin"
)

func (a Api) Operator(c *gin.Context) {
	agentID := c.Query("agent-id")
	if agentID == "" {
		helper.FailedWithJson(http.StatusBadRequest,
			exception.New(exceptionCode.ErrorBadRequest, exception.WithMsg("agent-id is empty")), c)
		return
	}
	gatewayID := c.Query("gateway-id")
	if gatewayID == "" {
		helper.FailedWithJson(http.StatusBadRequest,
			exception.New(exceptionCode.ErrorBadRequest, exception.WithMsg("gateway-id is empty")), c)
		return
	}

	claims, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.New(exceptionCode.ErrorGetClaimsFailed, exception.WithMsg("claims is empty")), c)
		return
	}
	userid := claims.User.Id

	var (
		taskid string
		resp   response.AgentAsyncExecuteOperatorResponse
	)

	switch a._getOperator(c) {
	case unknownOperator:
		helper.FailedWithJson(http.StatusBadRequest,
			exception.New(exceptionCode.ErrorBadRequest, exception.WithMsg("operator is empty")), c)
		return
	case startAgent:
		taskid, err = a.IAgentService.StartAgent(c, userid, gatewayID, agentID)
		if err != nil {
			helper.FailedWithJson(http.StatusInternalServerError,
				exception.NewWithErr(err, exceptionCode.StartAgentFailed), c)
			return
		}
	case stopAgent:
		taskid, err = a.IAgentService.StopAgent(c, userid, gatewayID, agentID)
		if err != nil {
			helper.FailedWithJson(http.StatusInternalServerError,
				exception.NewWithErr(err, exceptionCode.StopAgentFailed), c)
			return
		}
	case restartAgent:
		// todo
	case startGather:
		taskid, err = a.IAgentService.StartGather(c, userid, gatewayID, agentID)
		if err != nil {
			helper.FailedWithJson(http.StatusInternalServerError,
				exception.NewWithErr(err, exceptionCode.StartAgentFailed, exception.WithMsg("failed to start gather")), c)
			return
		}
	case stopGather:
		taskid, err = a.IAgentService.StopGather(c, userid, gatewayID, agentID)
		if err != nil {
			helper.FailedWithJson(http.StatusInternalServerError,
				exception.NewWithErr(err, exceptionCode.StopAgentFailed), c)
			return
		}
	case getOpenapiDoc:
		doc, err := a.IAgentService.GetOpenApiDoc(c, userid, gatewayID, agentID, 0)
		if err != nil {
			helper.FailedWithJson(http.StatusInternalServerError,
				exception.NewWithErr(err, exceptionCode.GetOpenAPIDocFailed), c)
			return
		}
		helper.SuccessJson(doc, c)
	}
	resp.TaskID = taskid
	helper.SuccessJson(resp, c)
}
func (a Api) GetAgentInfo(c *gin.Context) {
	agentID := c.Query("agent-id")
	if agentID == "" {
		helper.FailedWithJson(http.StatusBadRequest,
			exception.New(exceptionCode.ErrorBadRequest, exception.WithMsg("agent-id is empty")), c)
		return
	}
	gatewayID := c.Query("gateway-id")
	if gatewayID == "" {
		helper.FailedWithJson(http.StatusBadRequest,
			exception.New(exceptionCode.ErrorBadRequest, exception.WithMsg("gateway-id is empty")), c)
		return
	}

	claims, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.New(exceptionCode.ErrorGetClaimsFailed, exception.WithMsg("claims is empty")), c)
		return
	}
	userid := claims.User.Id

	agentInfo, err := a.IAgentService.GetAgentInfo(c, userid, gatewayID, agentID, true)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.NewWithErr(err, exceptionCode.GetAgentInfoFailed), c)
		return
	}
	helper.SuccessJson(agentInfo, c)
}
