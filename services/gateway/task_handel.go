package gateway

import (
	"fmt"

	"github.com/aenjoy/iot-lubricant/pkg/cache"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"

	"github.com/bytedance/sonic"
	grpcStatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (a *app) handelTask(task *corepb.TaskDetail, c *cache.MemoryCache[*corepb.QueryTaskResultResponse]) {
	logg.L.Debugf("running task ID:%s Message:%s Type:%v", task.TaskId, task.MessageId, task.GetTask())
	working := new(corepb.QueryTaskResultResponse_Working)
	finish := new(corepb.QueryTaskResultResponse_Finish)
	failed := new(corepb.QueryTaskResultResponse_Failed)
	working.Working = new(grpcStatus.Status)
	var result = &corepb.QueryTaskResultResponse{
		TaskId: task.TaskId,
		Result: working,
	}
	defer func() {
		_reportMessage <- &corepb.ReportRequest{
			Req: &corepb.ReportRequest_TaskResult{
				TaskResult: &corepb.TaskResultRequest{
					Msg: result,
				},
			},
		}
	}()

	c.Set(task.GetTaskId(), task.GetTaskId(), cache.NewStoreResult(cache.NeverExpired, result))

	setWorkingStatus := func(status string) {
		wor, _ := anypb.New(wrapperspb.String(status))
		working.Working.Details = []*anypb.Any{wor}
	}
	switch t := task.GetTask().(type) {
	case *corepb.TaskDetail_StartAgentRequest:
		logg.L.Debugf("StartAgentRequest")
		ids := t.StartAgentRequest.GetAgentId()
		working.Working.Details = make([]*anypb.Any, len(ids))
		for i := 0; i < len(ids); i++ {
			working.Working.Details[i], _ = anypb.New(wrapperspb.String("pending"))
		}

		// todo:这里可以优化为并发执行
		for i, id := range ids {
			working.Working.Details[i], _ = anypb.New(wrapperspb.String("working"))
			err := a.agent.StartAgent(id)
			if err != nil {
				working.Working.Details[i], _ = anypb.New(wrapperspb.String(fmt.Sprintf("failed due to:%v", err)))
				failed.Failed = &grpcStatus.Status{
					//Code:    int32(err.Code()),
					Message: err.Error(),
				}
				result.Result = failed
			} else {
				working.Working.Details[i], _ = anypb.New(wrapperspb.String("done"))
			}
		}
	case *corepb.TaskDetail_CreateAgentRequest:
		logg.L.Debugf("CreateAgentRequest")
		setWorkingStatus("creating")
		req := model.ProxypbCreateAgentRequest2CreateAgentRequest(t.CreateAgentRequest)
		err := a.agent.CreateAgent(req)
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			logg.L.Errorf("failed to create ot add agent: %v\n", err)
			result.Result = failed
			return
		}
		setWorkingStatus("done")
	case *corepb.TaskDetail_EditAgentRequest:
		logg.L.Debugf("EditAgentRequest")
		setWorkingStatus("editing")
		err := a.agent.EditAgent(t.EditAgentRequest.GetAgentId(), t.EditAgentRequest)
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			logg.L.Errorf("failed to edit agent: %v\n", err)
			result.Result = failed
			return
		}
		setWorkingStatus("done")
	case *corepb.TaskDetail_RemoveAgentRequest:
		logg.L.Debugf("RemoveAgentRequest")
		// todo:这里可以优化为并发执行
		ids := t.RemoveAgentRequest.GetAgentId()
		working.Working.Details = make([]*anypb.Any, len(ids))
		for i := 0; i < len(ids); i++ {
			working.Working.Details[i], _ = anypb.New(wrapperspb.String("pending"))
		}
		for i, id := range ids {
			working.Working.Details[i], _ = anypb.New(wrapperspb.String("removing"))
			err := a.agent.RemoveAgent(id)
			if err != nil {
				working.Working.Details[i], _ = anypb.New(wrapperspb.String(fmt.Sprintf("failed due to:%v", err)))
				failed.Failed = &grpcStatus.Status{
					//Code:    int32(err.Code()),
					Message: err.Error(),
				}
				result.Result = failed
			} else {
				working.Working.Details[i], _ = anypb.New(wrapperspb.String("done"))
			}
		}
	case *corepb.TaskDetail_StopAgentRequest:
		logg.L.Debugf("StopAgentRequest")
		ids := t.StopAgentRequest.GetAgentId()
		working.Working.Details = make([]*anypb.Any, len(ids))
		for i := 0; i < len(ids); i++ {
			working.Working.Details[i], _ = anypb.New(wrapperspb.String("pending"))
		}
		for i, id := range ids {
			working.Working.Details[i], _ = anypb.New(wrapperspb.String("stopping"))
			err := a.agent.StopAgent(id)
			if err != nil {
				working.Working.Details[i], _ = anypb.New(wrapperspb.String(fmt.Sprintf("failed due to:%v", err)))
				failed.Failed = &grpcStatus.Status{
					//Code:    int32(err.Code()),
					Message: err.Error(),
				}
				result.Result = failed
			} else {
				working.Working.Details[i], _ = anypb.New(wrapperspb.String("done"))
			}
		}
	case *corepb.TaskDetail_UpdateAgentRequest:
		logg.L.Debugf("UpdateAgentRequest")
		setWorkingStatus("updating")
		var conf *model.CreateAgentRequest
		if data := t.UpdateAgentRequest.GetConf(); len(data) > 0 {
			conf = &model.CreateAgentRequest{CreateAgentConf: new(model.CreateAgentConf)}
			err := sonic.Unmarshal(data, conf.CreateAgentConf)
			if err != nil {
				logg.L.Errorf("failed to unmarshal the conf:%v\n", err)
				return
			}
		}
		err := a.agent.UpdateAgent(t.UpdateAgentRequest.GetAgentId(), conf)
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			return
		}
		setWorkingStatus("done")
	case *corepb.TaskDetail_GetAgentStatusRequest:
		logg.L.Debugf("GetAgentStatusRequest")
		ids := t.GetAgentStatusRequest.GetAgentId()
		working.Working.Details = make([]*anypb.Any, len(ids))
		for i, id := range ids {
			working.Working.Details[i], _ = anypb.New(
				wrapperspb.String(a.agent.GetAgentStatus(id).String()))
		}
	case *corepb.TaskDetail_StartGatherRequest:
		logg.L.Debugf("StartGatherRequest")
		setWorkingStatus("starting")
		err := a.agent.StartGather(t.StartGatherRequest.GetAgentId())
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			result.Result = failed
			return
		}
		setWorkingStatus("done")
	case *corepb.TaskDetail_StopGatherRequest:
		logg.L.Debugf("StopGatherRequest")
		setWorkingStatus("stopping")
		err := a.agent.StopGather(t.StopGatherRequest.GetAgentId())
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			result.Result = failed
			return
		}
		setWorkingStatus("done")
	case *corepb.TaskDetail_GetAgentOpenAPIDocRequest:
		logg.L.Debugf("GetAgentOpenAPIDocRequest")
		setWorkingStatus("getting")
		doc, err := a.agent.GetAgentOpenApiDoc(t.GetAgentOpenAPIDocRequest.GetReq())
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			result.Result = failed
			return
		}

		a, _ := anypb.New(doc)
		working.Working.Details = []*anypb.Any{a}
	case *corepb.TaskDetail_GetAgentInfoRequest:
		logg.L.Debugf("GetAgentInfoRequest")
		setWorkingStatus("getting")
		info, err := a.agent.GetAgentInfo(t.GetAgentInfoRequest.GetAgentId())
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			result.Result = failed
			return
		}

		a, _ := anypb.New(info)
		working.Working.Details = []*anypb.Any{a}
	case *corepb.TaskDetail_SetAgentInfoRequest:
		logg.L.Debugf("SetAgentInfoRequest")
		setWorkingStatus("setting")
		req := t.SetAgentInfoRequest.GetInfo()
		err := a.agent.SetAgent(req.GetAgentID(), req)
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			result.Result = failed
			return
		}
		setWorkingStatus("done")
	case *corepb.TaskDetail_GetAgentIsGatheringRequest:
		logg.L.Debugf("GetAgentIsGatheringRequest")
		setWorkingStatus("getting")
		gatheringStatusResult, err := a.agent.IsGathering(t.GetAgentIsGatheringRequest.GetAgentId())
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			result.Result = failed
			return
		}
		a, _ := anypb.New(wrapperspb.Bool(gatheringStatusResult))
		working.Working.Details = []*anypb.Any{a}
	default:
		logg.L.Errorf("upsupport task type: %v", t)
		setWorkingStatus(fmt.Sprintf("failed due to: upsupport task type: %v", t))
		result.Result = failed
		return
	}

	finish.Finish, _ = anypb.New(working.Working)
	result.Result = finish
}
