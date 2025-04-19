package helper

import (
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/gin-gonic/gin"
)

func RequestBind[T any](c *gin.Context) *T {
	var req T
	err := c.BindJSON(&req)
	if err != nil {
		FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, code.ErrorBind), c)
		return nil
	}
	return &req
}
