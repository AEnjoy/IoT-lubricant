package monitor

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/helper"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

func (a Api) BaseInfo(c *gin.Context) { // 这个要加缓存中间件
	claims, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedByServer(err, c)
		return
	}
	userid := claims.User.Id
	key := fmt.Sprintf("baseinfo-query-user-%s", userid)

	// dev stg, ignore cache
	//result, _ := a.DataStore.CacheCli.Get(c, key)
	//if result != "" {
	//	if helper.JsonString(http.StatusOK, result, "success", "0000", c) {
	//		return
	//	}
	//	// no cache or cache expired/failed
	//}

	var (
		gateways []model.Gateway

		wg     sync.WaitGroup
		output response.QueryMonitorBaseInfoResponse
	)

	gateways, err = a.DataStore.ICoreDb.GetAllGatewayByUserID(c, userid)
	if err != nil {
		e := exception.ErrNewException(err, exceptionCode.GetGatewayFailed)
		logger.Errorf("GetAllGatewayByUserID failed err: %v", e)
		helper.FailedWithJson(http.StatusInternalServerError, e, c)
		return
	}

	for _, gateway := range gateways {
		output.GatewayCount++
		if gateway.Status != "online" {
			output.OfflineGateway++
		}
		list, err := a.DataStore.ICoreDb.GetAgentList(c, userid, gateway.GatewayID)
		if err != nil {
			logger.Errorf("GetAgentList failed,err: %v", err)
			continue
		}

		var ids []string
		for _, agent := range list {
			ids = append(ids, agent.AgentId)
			output.AgentCount++
		}

		wg.Add(1)
		go func(c *gin.Context, id string, ids []string, output *response.QueryMonitorBaseInfoResponse) {
			defer wg.Done()
			status, _ := a.IAgentService.GetAgentStatus(c, id, ids)
			for _, agentStatus := range status {
				if agentStatus != model.StatusOnline {
					atomic.AddInt32(&output.OfflineAgent, 1)
				}
			}
		}(c, gateway.GatewayID, ids, &output)
	}

	wg.Wait()
	str, _ := sonic.MarshalString(&output)

	err = a.CacheCli.SetEx(c, key, str, 10*time.Minute)
	if err != nil {
		logger.Errorf("set cache failed,err: %v", err)
	}
	helper.SuccessJson(output, c)
}
