package task

import (
	"net/http"
	"strconv"

	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func (a Api) GetTaskList(c *gin.Context) {
	//	Current  int `json:"current"`  // 当前页
	//	PageSize int `json:"pageSize"` // 数据条数
	// todo
	claims, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.New(exceptionCode.ErrorGetClaimsFailed, exception.WithMsg("claims is empty")), c)
		return
	}

	userid := claims.User.Id
	_current := c.DefaultQuery("current", "1")
	_size := c.DefaultQuery("pageSize", "10")

	current, err := strconv.Atoi(_current)
	if err != nil {
		helper.FailedWithJson(http.StatusBadRequest, exception.ErrNewException(err,
			exceptionCode.ErrorBadRequest,
			exception.WithMsg("failed to parse `current`")),
			c)
		return
	}
	size, err := strconv.Atoi(_size)
	if err != nil {
		helper.FailedWithJson(http.StatusBadRequest, exception.ErrNewException(err,
			exceptionCode.ErrorBadRequest,
			exception.WithMsg("failed to parse `pageSize`")),
			c)
		return
	}
	resp, err := a.ICoreDb.UserGetAsyncJobs(c, userid, current, size)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError, exception.ErrNewException(err, exceptionCode.ErrorInternalServerError), c)
		return
	}
	helper.SuccessJson(resp, c)

}
func (a Api) QueryTask(c *gin.Context) {
	taskID := c.Query("taskId")
	if taskID == "" {
		helper.FailedWithJson(http.StatusBadRequest, exception.New(exceptionCode.ErrorBadRequest, exception.WithMsg("taskId is empty")), c)
		return
	}
	var resp response.QueryTaskResultResponse
	resp.TaskID = taskID
	status, result, err := a.ICoreDb.GetAsyncJobResult(c, taskID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.Result = "taskId not found"
			helper.SuccessJson(resp, c)
			return
		}
		helper.FailedWithJson(http.StatusInternalServerError, exception.ErrNewException(err, exceptionCode.ErrorInternalServerError), c)
		return
	}

	resp.Status = status
	resp.Result = result

	helper.SuccessJson(resp, c)
}
