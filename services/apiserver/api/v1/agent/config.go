package agent

import (
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/helper"

	"github.com/gin-gonic/gin"
)

func (a Api) Set(c *gin.Context) {
	agentInfo := helper.RequestBind[agentpb.AgentInfo](c)
	if agentInfo == nil {
		return
	}

	agentID := agentInfo.GetAgentID()
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

	taskid, err := a.IAgentService.SetAgentInfo(c, userid, gatewayID, agentID, agentInfo)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.ErrorSetAgentInfoFailed), c,
		)
		return
	}

	helper.SuccessJson(response.AgentAsyncExecuteOperatorResponse{TaskID: taskid}, c)
}
