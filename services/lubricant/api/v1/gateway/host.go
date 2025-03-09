package gateway

import (
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/ssh"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/gin-gonic/gin"
)

func (a Api) AddHost(c *gin.Context) {
	_, gatewayHostInfo := a.getGatewayHostModel(c)
	if gatewayHostInfo == nil {
		return
	}
	_, err := ssh.NewSSHClient(gatewayHostInfo, true)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err,
				exceptionCode.LinkToGatewayFailed,
				exception.WithMsg("test linker failed: failed to link to target host")), c)
		return
	}
	err = a.IGatewayService.AddHost(c, gatewayHostInfo)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.AddGatewayHostFailed), c)
		return
	}
	helper.SuccessJson(gatewayHostInfo, c)
}
func (a Api) ListHosts(c *gin.Context) {
	claimsUser, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedByServer(err, c)
		return
	}
	hosts, err := a.IGatewayService.UserGetHosts(c, claimsUser.User.Id)
	if err != nil {
		helper.FailedWithJson(http.StatusNotFound, exception.ErrNewException(err,
			exceptionCode.ErrorNotFound), c)
		return
	}
	helper.SuccessJson(hosts, c)
}

func (a Api) DescriptionHost(c *gin.Context) {
	hostid, ok := c.GetQuery("host_id")
	if !ok {
		helper.FailedWithJson(http.StatusBadRequest, exception.New(
			exceptionCode.ErrorBadRequest, exception.WithMsg("host_id is required")), c)
		return
	}
	host, err := a.IGatewayService.DescriptionHost(c, hostid)
	if err != nil {
		helper.FailedWithJson(http.StatusBadRequest, exception.ErrNewException(err,
			exceptionCode.DescriptionHostFailed, exception.WithData(host)), c)
		return
	}
	helper.SuccessJson(host, c)
}
