package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	gatewaypb "github.com/aenjoy/iot-lubricant/protobuf/gateway"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/rs/xid"
	"google.golang.org/genproto/googleapis/rpc/status"
)

var _true = true

//	func (a *AgentService) GetAgentInfoByTaskID(ctx context.Context, userid, gatewayID, agentID, id string) (*agentpb.AgentInfo, error) {
//		if id == "" {
//			return a.GetAgentInfo(ctx, userid, gatewayID, agentID, true)
//		}
//		a.store.ICoreDb.GetAsyncJobResult(ctx, id) // todo:need refact database store
//
// }
func (a *AgentService) IsGathering(ctx context.Context, userid, gatewayID, agentID string) (bool, error) {
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId:            id,
		IsSynchronousTask: &_true,
		Task: &corepb.TaskDetail_GetAgentIsGatheringRequest{
			GetAgentIsGatheringRequest: &gatewaypb.GetAgentIsGatheringRequest{
				AgentId: agentID,
			},
		},
	}
	_, _, err := a.PushTaskAgentPb(ctx, &id, userid, gatewayID, agentID, td)
	if err != nil {
		_ = a.store.ICoreDb.SetAsyncJobStatus(ctx, nil, id, "failed", "failed to push task")
		return false, err
	}
	resp, err := a.SyncTaskQueue.WaitTask(id, 10*time.Second)
	if err != nil {
		_ = a.store.ICoreDb.SetAsyncJobStatus(ctx, nil, id, "failed", "failed to wait task")
		return false, err
	}
	if resp.GetFinish() != nil {
		var s status.Status
		var isGathering wrapperspb.BoolValue

		if err := resp.GetFinish().UnmarshalTo(&s); err != nil {
			return false, err
		}
		if len(s.Details) == 0 {
			return false, fmt.Errorf("get agent info failed: %v", resp.GetResult())
		}

		if err := s.Details[0].UnmarshalTo(&isGathering); err != nil {
			return false, err
		}
		return isGathering.GetValue(), nil
	}
	str := fmt.Sprintf("get agent info failed: %v", resp.GetResult())
	_ = a.store.ICoreDb.SetAsyncJobStatus(ctx, nil, id, "failed", str)
	return false, errors.New(str)
}
func (a *AgentService) GetAgentInfo(ctx context.Context, userid string, gatewayID string, agentID string, sync bool) (*agentpb.AgentInfo, error) {
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId:            id,
		IsSynchronousTask: &sync,
		Task: &corepb.TaskDetail_GetAgentInfoRequest{
			GetAgentInfoRequest: &gatewaypb.GetAgentInfoRequest{
				AgentId: agentID,
			},
		},
	}
	_, _, err := a.PushTaskAgentPb(ctx, &id, userid, gatewayID, agentID, td)
	if err != nil {
		_ = a.store.ICoreDb.SetAsyncJobStatus(ctx, nil, id, "failed", "failed to push task")
		return nil, err
	}
	if !sync {
		return nil, nil
	}

	resp, err := a.SyncTaskQueue.WaitTask(id, 10*time.Second)
	if err != nil {
		_ = a.store.ICoreDb.SetAsyncJobStatus(ctx, nil, id, "failed", fmt.Sprintf("get agentInfo failed: %v", err))
		return nil, err
	}

	if resp.GetFinish() != nil {
		var s status.Status
		var info agentpb.AgentInfo

		if err := resp.GetFinish().UnmarshalTo(&s); err != nil {
			return nil, err
		}
		if len(s.Details) == 0 {
			return nil, fmt.Errorf("get agent info failed: %v", resp.GetResult())
		}

		if err := s.Details[0].UnmarshalTo(&info); err != nil {
			return nil, err
		}
		return &info, nil
	}
	str := fmt.Sprintf("get agent info failed: %v", resp.GetResult())
	_ = a.store.ICoreDb.SetAsyncJobStatus(ctx, nil, id, "failed", str)
	return nil, errors.New(str)
}
func (a *AgentService) ListAgents(ctx context.Context, userID, gatewayid string) ([]model.Agent, error) {
	return a.store.GetAgentList(ctx, userID, gatewayid)
}
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
func (a *AgentService) SetAgentInfo(ctx context.Context, userid, gatewayid, agentid string, info *agentpb.AgentInfo) (taskid string, err error) {
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId: id,
		Task: &corepb.TaskDetail_SetAgentInfoRequest{
			SetAgentInfoRequest: &gatewaypb.SetAgentInfoRequest{
				Info: info,
			},
		},
	}
	_, _, err = a.PushTaskAgentPb(ctx, &id, userid, gatewayid, agentid, td)
	return id, err
}
func (a *AgentService) StartAgent(ctx context.Context, userid, gatewayid, agentid string) (taskid string, err error) {
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId: id,
		Task: &corepb.TaskDetail_StartAgentRequest{
			StartAgentRequest: &gatewaypb.StartAgentRequest{
				AgentId: []string{agentid},
			},
		},
	}
	_, _, err = a.PushTaskAgentPb(ctx, &id, userid, gatewayid, agentid, td)
	return id, err
}

func (a *AgentService) StopAgent(ctx context.Context, userid, gatewayid, agentid string) (taskid string, err error) {
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId: id,
		Task: &corepb.TaskDetail_StopAgentRequest{
			StopAgentRequest: &gatewaypb.StopAgentRequest{
				AgentId: []string{agentid},
			},
		},
	}
	_, _, err = a.PushTaskAgentPb(ctx, &id, userid, gatewayid, agentid, td)
	return id, err
}

func (a *AgentService) StartGather(ctx context.Context, userid, gatewayid, agentid string) (taskid string, err error) {
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId: id,
		Task: &corepb.TaskDetail_StartGatherRequest{
			StartGatherRequest: &gatewaypb.StartGatherRequest{
				AgentId: agentid,
			},
		},
	}
	_, _, err = a.PushTaskAgentPb(ctx, &id, userid, gatewayid, agentid, td)
	return id, err
}

func (a *AgentService) StopGather(ctx context.Context, userid, gatewayid, agentid string) (taskid string, err error) {
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId: id,
		Task: &corepb.TaskDetail_StopGatherRequest{
			StopGatherRequest: &gatewaypb.StopGatherRequest{
				AgentId: agentid,
			},
		},
	}
	_, _, err = a.PushTaskAgentPb(ctx, &id, userid, gatewayid, agentid, td)
	return id, err
}

func (a *AgentService) GetOpenApiDoc(ctx context.Context, userid, gatewayid, agentid string, docType agentpb.OpenapiDocType) (result *response.GetOpenApiDocResponse, err error) {
	id := xid.New().String()
	td := &corepb.TaskDetail{
		TaskId:            id,
		IsSynchronousTask: &_true,
		Task: &corepb.TaskDetail_GetAgentOpenAPIDocRequest{
			GetAgentOpenAPIDocRequest: &gatewaypb.GetAgentOpenAPIDocRequest{
				Req: &agentpb.GetOpenapiDocRequest{
					AgentID: agentid,
					DocType: docType,
				},
			},
		},
	}

	_, _, err = a.PushTaskAgentPb(ctx, &id, userid, gatewayid, agentid, td)
	if err != nil {
		return nil, err
	}

	resp, err := a.SyncTaskQueue.WaitTask(id, 10*time.Second)
	if err != nil {
		_ = a.store.ICoreDb.SetAsyncJobStatus(ctx, nil, id, "failed", fmt.Sprintf("get openapi doc failed: %v", err))
		return nil, err
	}

	if resp.GetFinish() != nil {
		var s status.Status
		var doc agentpb.OpenapiDoc

		if err := resp.GetFinish().UnmarshalTo(&s); err != nil {
			return nil, err
		}
		if len(s.Details) == 0 {
			return nil, fmt.Errorf("get openapi doc failed: %v", resp.GetResult())
		}

		if err := s.Details[0].UnmarshalTo(&doc); err != nil {
			return nil, err
		}
		return &response.GetOpenApiDocResponse{
			AgentID: agentid,
			Doc:     doc.GetOriginalFile(),
		}, nil
	}
	str := fmt.Sprintf("get openapi doc failed: %v", resp.GetResult())
	_ = a.store.ICoreDb.SetAsyncJobStatus(ctx, nil, id, "failed", str)
	return nil, errors.New(str)
}
