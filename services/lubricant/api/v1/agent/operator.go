package agent

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/bytedance/sonic"
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

	key := fmt.Sprintf("getAgentInfo-user-%s-agent-%s-gateway-%s", userid, agentID, gatewayID)
	//dev stg, ignore cache
	result, _ := a.DataStore.CacheCli.Get(c, key)
	if result != "" {
		if helper.JsonString(http.StatusOK, result, "success", "0000", c) {
			return
		}
		// no cache or cache expired/failed
	}
	agentInfo, err := a.IAgentService.GetAgentInfo(c, userid, gatewayID, agentID, true)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.NewWithErr(err, exceptionCode.GetAgentInfoFailed), c)
		return
	}

	str, _ := sonic.MarshalString(&agentInfo)

	err = a.CacheCli.SetEx(c, key, str, 10*time.Minute)
	if err != nil {
		logger.Errorf("set cache failed,err: %v", err)
	}
	helper.JsonString(http.StatusOK, str, "success", "0000", c)
}
func (a Api) List(c *gin.Context) {
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
	//var resp response.ListAgentResponse
	agents, err := a.IAgentService.ListAgents(c, userid, gatewayID)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.NewWithErr(err, exceptionCode.ListAgentFailed), c)
		return
	}
	helper.SuccessJson(response.ListAgentResponse{Agents: agents, GatewayID: gatewayID}, c)
}
