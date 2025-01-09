package gateway

import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/internal/cache"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (a *app) handelTask(task *core.TaskDetail, c *cache.MemoryCache[*core.QueryTaskResultResponse]) {
	// todo:impl me
	working := new(core.QueryTaskResultResponse_Working)
	finish := new(core.QueryTaskResultResponse_Finish)
	working.Working = new(status.Status)
	var result = &core.QueryTaskResultResponse{
		TaskId: task.TaskId,
		Result: working,
	}
	c.Set(task.GetTaskId(), "", cache.NewStoreResult(cache.NeverExpired, result))

	switch t := task.GetTask().(type) {
	case *core.TaskDetail_StartAgentRequest:
		ids := t.StartAgentRequest.GetAgentId()
		working.Working.Details = make([]*anypb.Any, len(ids), len(ids))
		for i := 0; i < len(ids); i++ {
			working.Working.Details[i], _ = anypb.New(wrapperspb.String("pending"))
		}

		// todo:这里可以优化为并发执行
		for i, id := range ids {
			working.Working.Details[i], _ = anypb.New(wrapperspb.String("working"))
			err := a.agent.StartAgent(id)
			if err != nil {
				working.Working.Details[i], _ = anypb.New(wrapperspb.String(fmt.Sprintf("failed due to:%v", err)))
			} else {
				working.Working.Details[i], _ = anypb.New(wrapperspb.String("done"))
			}
		}

		finish.Finish, _ = anypb.New(working.Working)
		result.Result = finish
	case *core.TaskDetail_CreateAgentRequest:
		setStatus := func(status string) {
			wor, _ := anypb.New(wrapperspb.String(status))
			working.Working.Details = []*anypb.Any{wor}
		}
		setStatus("creating")
		req := model.ProxypbCreateAgentRequest2CreateAgentRequest(t.CreateAgentRequest)
		a.agent.CreateAgent(req)
	case *core.TaskDetail_EditAgentRequest:
	case *core.TaskDetail_RemoveAgentRequest:
	case *core.TaskDetail_StopAgentRequest:
	case *core.TaskDetail_UpdateAgentRequest:
	}
}
