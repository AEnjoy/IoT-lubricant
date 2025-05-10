package gateway

import (
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/model/request"
	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/helper"

	"github.com/cloudwego/base64x"
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
	task, err := base64x.StdEncoding.DecodeString(req.Task)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.ErrorDecodeFailed,
				exception.WithMsg("error in task file")), c)
		return
	}
	claims, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.New(exceptionCode.ErrorGetClaimsFailed, exception.WithMsg("claims is empty")), c)
		return
	}
	userid := claims.User.Id

	_, _, err = a.IAgentService.PushTaskAgent(c, taskid, userid, gatewayID, req.AgentID, task)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.ErrorPushTaskFailed), c,
		)
		return
	}
	helper.SuccessJson(response.PushAgentTaskResponse{TaskID: *taskid}, c)
}
