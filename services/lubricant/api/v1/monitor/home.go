package monitor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

func (a Api) BaseInfo(c *gin.Context) { // 这个要加缓存中间件
	var output response.QueryMonitorBaseInfoResponse

	claims, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedByServer(err, c)
		return
	}

	key := fmt.Sprintf("baseinfo-query-user-%s", claims.User.Id)
	result, _ := a.DataStore.CacheCli.Get(c, key)
	if result != "" {
		if helper.JsonString(http.StatusOK, result, "success", "0000", c) {
			return
		}
		// no cache or cache expired/failed
	}

	var (
		gateways []model.Gateway
		agents   []model.Agent
	)

	gateways, err = a.DataStore.ICoreDb.GetAllGatewayByUserID(c, claims.User.Id)
	if err != nil {
		e := exception.ErrNewException(err, exceptionCode.GetGatewayFailed)
		logger.Errorf("GetAllGatewayByUserID failed err: %v", e)
		helper.FailedWithJson(http.StatusInternalServerError, e, c)
		return
	}

	for _, gateway := range gateways {
		if gateway.Status != "running" {
			output.OfflineGateway++
		}
		list, err := a.DataStore.ICoreDb.GetAgentList(c, gateway.GatewayID)
		if err != nil {
			logger.Errorf("GetAgentList failed,err: %v", err)
			continue
		}
		agents = append(agents, list...)
	}
	output.GatewayCount = len(gateways)
	output.AgentCount = len(agents)
	// todo: agent 离线数据

	str, _ := sonic.MarshalString(&output)

	err = a.CacheCli.SetEx(c, key, str, 10*time.Minute)
	if err != nil {
		logger.Errorf("set cache failed,err: %v", err)
	}
	helper.SuccessJson(output, c)
}
