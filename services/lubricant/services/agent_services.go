package services

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
)

func (a *AgentService) GetAgentStatus(ctx context.Context, gatewayid string, ids []string) ([]model.AgentStatus, error) {
	gatewayStatus, err := a.db.GetGatewayStatus(ctx, gatewayid)
	if err != nil {
		err = exception.NewWithErr(err, exceptionCode.ErrorGetGatewayStatusFailed,
			exception.WithMsg("Failed to get gateway status"))
		return nil, err
	}

	var retVal = make([]model.AgentStatus, len(ids))
	if gatewayStatus != model.StatusOnline.String() {
		for i := range retVal {
			retVal[i] = model.StatusOffline
		}
		return retVal, nil
	}

	for i, id := range ids {
		agentStatus, err := a.db.GetAgentStatus(ctx, id)
		if err != nil {
			retVal[i] = model.StatusUnknown
			continue
		}
		retVal[i] = model.AgentStatus(agentStatus)
	}

	return retVal, nil
}

func (a *AgentService) StartAgent(ctx context.Context, gatewayid, agentid string) (taskid string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *AgentService) StopAgent(ctx context.Context, gatewayid, agentid string) (taskid string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *AgentService) StartGather(ctx context.Context, gatewayid, agentid string) (taskid string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *AgentService) StopGather(ctx context.Context, gatewayid, agentid string) (taskid string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *AgentService) GetOpenApiDoc(ctx context.Context, gatewayid, agentid string) (result string, err error) {
	//TODO implement me
	panic("implement me")
}
