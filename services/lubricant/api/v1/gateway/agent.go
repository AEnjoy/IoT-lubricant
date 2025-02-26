package gateway

import (
	"encoding/base64"
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/form/request"
	"github.com/aenjoy/iot-lubricant/pkg/form/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/gin-gonic/gin"
)

func (a Api) AgentPushTask(c *gin.Context) {
	req := helper.RequestBind[request.PushTaskRequest](c)
	gatewayID := c.Param("gatewayid")
	taskid := func() *string {
		if req.TaskID == "" {
			return nil
		}
		return &req.TaskID
	}()
	task, err := base64.StdEncoding.DecodeString(req.Task)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.ErrorDecodeFailed,
				exception.WithMsg("error in task file")), c)
		return
	}
	_, _, err = a.IAgentService.PushTaskAgent(c, taskid, gatewayID, req.AgentID, task)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.ErrorPushTaskFailed), c,
		)
		return
	}
	helper.SuccessJson(response.PushAgentTaskResponse{TaskID: *taskid}, c)
}
