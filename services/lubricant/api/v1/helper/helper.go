package helper

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	"github.com/bytedance/sonic"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
)

// JsonString data is json string failed return false and nothing to do, else return true
func JsonString(httpCode int, data, msg, code string, c *gin.Context) bool {
	var m map[string]any
	err := sonic.Unmarshal([]byte(data), &m)
	if err != nil {
		// FailedWithJson(http.StatusInternalServerError, exception.ErrNewException(err,exceptionCode.ErrorDecodeJSON),c)
		return false
	}
	c.JSON(httpCode, response.Meta{
		Code: code,
		Msg:  msg,
		Data: m,
	})
	return true
}
func SuccessJson(data any, c *gin.Context) {
	c.JSON(http.StatusOK, response.Success{
		Meta: response.Meta{
			Code: "0000",
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
			Code: strconv.Itoa(code),
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
