package main

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"

	grpcStatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func handleTaskResp(cli corepb.CoreServiceClient, ctx context.Context) {
	go func() {
		taskStream, err := cli.GetTask(ctx)
		if err != nil {
			logger.Warnf("Faild to create getTask stream: %v", err)
			return
		}
		for {
			select {
			case <-ctx.Done():
				logger.Info("GetTask stream closed")
				return
			default:
			}
			recv, err := taskStream.Recv()
			if err != nil {
				logger.Warnf("Failed to receive task response: %v", err)
				continue
			}
			if req := recv.GetCorePushTaskRequest(); req != nil {
				executorTask(cli, ctx, req.GetMessage())
			}
		}
	}()
}
func executorTask(cli corepb.CoreServiceClient, ctx context.Context, task *corepb.TaskDetail) {
	working := new(corepb.QueryTaskResultResponse_Working)
	finish := new(corepb.QueryTaskResultResponse_Finish)
	//failed := new(corepb.QueryTaskResultResponse_Failed)
	working.Working = new(grpcStatus.Status)
	var result = &corepb.QueryTaskResultResponse{
		TaskId: task.TaskId,
		Result: working,
	}
	defer func() {
		finish.Finish, _ = anypb.New(working.Working)
		result.Result = finish

		req := &corepb.ReportRequest{
			Req: &corepb.ReportRequest_TaskResult{
				TaskResult: &corepb.TaskResultRequest{
					Msg: result,
				},
			},
			GatewayId: gatewayID,
		}
		// send report request
		_, err := cli.Report(ctx, req)
		if err != nil {
			logger.Errorf("Failed to send report request to server: %v", err)
		}
	}()
	setWorkingStatus := func(status string) {
		wor, _ := anypb.New(wrapperspb.String(status))
		working.Working.Details = []*anypb.Any{wor}
	}
	switch task.GetTask().(type) {
	case
			*corepb.TaskDetail_StartAgentRequest,
			*corepb.TaskDetail_CreateAgentRequest,
			*corepb.TaskDetail_EditAgentRequest,
			*corepb.TaskDetail_RemoveAgentRequest,
			*corepb.TaskDetail_StopAgentRequest,
			*corepb.TaskDetail_UpdateAgentRequest,
			*corepb.TaskDetail_GetAgentStatusRequest:
		setWorkingStatus("ok")
		a, _ := anypb.New(wrapperspb.String("ok"))
		working.Working.Details = []*anypb.Any{a}
	case *corepb.TaskDetail_GetAgentIsGatheringRequest:
		setWorkingStatus("ok")
		a, _ := anypb.New(wrapperspb.Bool(true))
		working.Working.Details = []*anypb.Any{a}
	case *corepb.TaskDetail_GetAgentInfoRequest:
		setWorkingStatus("ok")
		a, _ := anypb.New(&agentpb.AgentInfo{
			AgentID:   randGetAgentID(),
			GatewayID: &gatewayID,
		})
		working.Working.Details = []*anypb.Any{a}
	case *corepb.TaskDetail_GetAgentOpenAPIDocRequest:
		setWorkingStatus("ok")
		a, _ := anypb.New(&agentpb.OpenapiDoc{})
		working.Working.Details = []*anypb.Any{a}
	default:
	}
}
