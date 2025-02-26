package gateway

import (
	"encoding/base64"
	"net/http"

	"github.com/AEnjoy/IoT-lubricant/pkg/form/request"
	"github.com/AEnjoy/IoT-lubricant/pkg/form/response"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	helper2 "github.com/AEnjoy/IoT-lubricant/services/core/api/v1/helper"
	"github.com/gin-gonic/gin"
)

func (a Api) AgentPushTask(c *gin.Context) {
	req := helper2.RequestBind[request.PushTaskRequest](c)
	gatewayID := c.Param("gatewayid")
	taskid := func() *string {
		if req.TaskID == "" {
			return nil
		}
		return &req.TaskID
	}()
	task, err := base64.StdEncoding.DecodeString(req.Task)
	if err != nil {
		helper2.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.ErrorDecodeFailed,
				exception.WithMsg("error in task file")), c)
		return
	}
	_, _, err = a.IAgentService.PushTaskAgent(c, taskid, gatewayID, req.AgentID, task)
	if err != nil {
		helper2.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptCode.ErrorPushTaskFailed), c,
		)
		return
	}
	helper2.SuccessJson(response.PushAgentTaskResponse{TaskID: *taskid}, c)
}
