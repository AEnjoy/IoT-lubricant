package services

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
)

func (a *AgentService) GetAgentStatus(ctx context.Context, gatewayid string, ids []string) ([]model.AgentStatus, error) {
	status, err := a.db.GetGatewayStatus(ctx, gatewayid)
	if err != nil {
		err = exception.NewWithErr(err, exceptionCode.ErrorGetGatewayStatusFailed,
			exception.WithMsg("Failed to get gateway status"))
		return nil, err
	}

	var retVal = make([]model.AgentStatus, len(ids))
	if status != model.StatusOnline.String() {
		for i, _ := range retVal {
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
