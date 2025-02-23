package gateway

import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/internal/cache"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	corepb "github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/bytedance/sonic"
	grpcStatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (a *app) handelTask(task *corepb.TaskDetail, c *cache.MemoryCache[*corepb.QueryTaskResultResponse]) {
	logger.Debugf("running task ID:%s Message:%s Type:%v", task.TaskId, task.MessageId, task.GetTask())
	// todo:impl me
	working := new(corepb.QueryTaskResultResponse_Working)
	finish := new(corepb.QueryTaskResultResponse_Finish)
	failed := new(corepb.QueryTaskResultResponse_Failed)
	working.Working = new(grpcStatus.Status)
	var result = &corepb.QueryTaskResultResponse{
		TaskId: task.TaskId,
		Result: working,
	}
	c.Set(task.GetTaskId(), "", cache.NewStoreResult(cache.NeverExpired, result))

	setWorkingStatus := func(status string) {
		wor, _ := anypb.New(wrapperspb.String(status))
		working.Working.Details = []*anypb.Any{wor}
	}
	switch t := task.GetTask().(type) {
	case *corepb.TaskDetail_StartAgentRequest:
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
		setWorkingStatus("creating")
		req := model.ProxypbCreateAgentRequest2CreateAgentRequest(t.CreateAgentRequest)
		err := a.agent.CreateAgent(req)
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			logger.Errorf("failed to create ot add agent: %v\n", err)
			result.Result = failed
			return
		}
		setWorkingStatus("done")
	case *corepb.TaskDetail_EditAgentRequest:
		setWorkingStatus("editing")
		err := a.agent.EditAgent(t.EditAgentRequest.GetAgentId(), t.EditAgentRequest)
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			logger.Errorf("failed to edit agent: %v\n", err)
			result.Result = failed
			return
		}
		setWorkingStatus("done")
	case *corepb.TaskDetail_RemoveAgentRequest:
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
		setWorkingStatus("updating")
		var conf *model.CreateAgentRequest
		if data := t.UpdateAgentRequest.GetConf(); len(data) > 0 {
			conf = &model.CreateAgentRequest{CreateAgentConf: new(model.CreateAgentConf)}
			err := sonic.Unmarshal(data, conf.CreateAgentConf)
			if err != nil {
				logger.Errorf("failed to unmarshal the conf:%v\n", err)
				return
			}
		}
		err := a.agent.UpdateAgent(t.UpdateAgentRequest.GetAgentId(), conf)
		if err != nil {
			setWorkingStatus(fmt.Sprintf("failed due to:%v", err))
			return
		}
		setWorkingStatus("done")
	}

	finish.Finish, _ = anypb.New(working.Working)
	result.Result = finish
}
