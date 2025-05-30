package agent

import (
	"fmt"
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/helper"

	"github.com/gin-gonic/gin"
)

func (a Api) GetData(context *gin.Context) {
	agentId := context.Query("agent_id")
	if agentId == "" {
		helper.FailedWithJson(http.StatusBadRequest,
			exception.New(exceptionCode.ErrorNeedAgentID), context)
		return
	}

	project, err := a.DataStore.GetProjectByAgentID(context, agentId)
	if err != nil {
		helper.FailedWithErrorJson(http.StatusInternalServerError, err, context)
	}
	if project.ProjectID == "" {
		helper.FailedWithJson(http.StatusNotFound,
			exception.New(exceptionCode.ErrorGatewayAgentNotFound), context)
		return
	}

	data, _ := a.DataStore.CacheCli.Get(context, fmt.Sprintf(constant.LatestDataCacheKey, project.ProjectID, agentId))
	helper.SuccessJson(data, context)
}
