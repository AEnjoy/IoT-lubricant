package services

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	gatewaypb "github.com/aenjoy/iot-lubricant/protobuf/gateway"

	"github.com/rs/xid"
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
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId: id,
		Task: &corepb.TaskDetail_StartAgentRequest{
			StartAgentRequest: &gatewaypb.StartAgentRequest{
				AgentId: []string{agentid},
			},
		},
	}
	_, _, err = a.PushTaskAgentPb(ctx, &id, gatewayid, agentid, td)
	return id, err
}

func (a *AgentService) StopAgent(ctx context.Context, gatewayid, agentid string) (taskid string, err error) {
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId: id,
		Task: &corepb.TaskDetail_StopAgentRequest{
			StopAgentRequest: &gatewaypb.StopAgentRequest{
				AgentId: []string{agentid},
			},
		},
	}
	_, _, err = a.PushTaskAgentPb(ctx, &id, gatewayid, agentid, td)
	return id, err
}

func (a *AgentService) StartGather(ctx context.Context, gatewayid, agentid string) (taskid string, err error) {
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId: id,
		Task: &corepb.TaskDetail_StartGatherRequest{
			StartGatherRequest: &gatewaypb.StartGatherRequest{
				AgentId: agentid,
			},
		},
	}
	_, _, err = a.PushTaskAgentPb(ctx, &id, gatewayid, agentid, td)
	return id, err
}

func (a *AgentService) StopGather(ctx context.Context, gatewayid, agentid string) (taskid string, err error) {
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId: id,
		Task: &corepb.TaskDetail_StopGatherRequest{
			StopGatherRequest: &gatewaypb.StopGatherRequest{
				AgentId: agentid,
			},
		},
	}
	_, _, err = a.PushTaskAgentPb(ctx, &id, gatewayid, agentid, td)
	return id, err
}

func (a *AgentService) GetOpenApiDoc(ctx context.Context, gatewayid, agentid string) (result string, err error) {
	//TODO implement me
	panic("implement me")
}
