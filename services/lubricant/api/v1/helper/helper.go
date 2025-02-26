package helper

import (
	"errors"
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/form/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
)

func SuccessJson(data any, c *gin.Context) {
	c.JSON(http.StatusOK, response.Success{
		Meta: response.Meta{
			Msg:  "success",
			Data: data,
		},
	})
}
func SuccessWithData(data []byte, c *gin.Context) {
	c.Data(http.StatusOK, "application/octet-stream", data)
}
func SuccessWithText(data string, c *gin.Context) {
	c.String(http.StatusOK, data)
}
func FailedByServer(err error, c *gin.Context) {
	// 500 状态, 接口报错, 返回内容: Exception对象
	httpCode := http.StatusInternalServerError
	_failed(httpCode, err, c)
}
func FailedByClient(err error, c *gin.Context) {
	// 400 状态, 接口报错, 返回内容: ApiException对象
	httpCode := http.StatusBadRequest
	_failed(httpCode, err, c)
}
func _failed(httpCode int, err error, c *gin.Context) {
	var v *exception.Exception
	if errors.As(err, &v) {
		if v.Code != 0 {
			httpCode = int(v.Code)
		}
		c.JSON(httpCode, v)
	} else {
		c.JSON(httpCode, err)
	}
	c.Abort()
}
func FailedWithErrorJson(code int, err error, c *gin.Context) {
	c.JSON(code, err)
}
func FailedWithJson(code int, exception *exception.Exception, c *gin.Context) {
	c.JSON(code, response.Failed{
		Meta: response.Meta{
			Code: code,
			Msg:  "failed",
			Data: exception,
		},
	})
}
func FailedWithError(code int, err error, c *gin.Context) {
	c.String(code, err.Error())
}
func GetClaims(c *gin.Context) (*casdoorsdk.Claims, error) {
	v, ok := c.Get("claims")
	if !ok {
		return nil, errors.New("claims not found")
	}
	return v.(*casdoorsdk.Claims), nil
}
