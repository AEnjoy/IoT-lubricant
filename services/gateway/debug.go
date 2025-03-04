package gateway

import (
	"context"
	"net"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/grpc/middleware"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	gatewaypb "github.com/aenjoy/iot-lubricant/protobuf/gateway"
	"github.com/aenjoy/iot-lubricant/services/gateway/services/agent"
	"github.com/aenjoy/iot-lubricant/services/gateway/services/async"
	"github.com/aenjoy/iot-lubricant/services/gateway/services/data"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

type gatewayDebugServer struct {
	gatewaypb.UnimplementedDebugServiceServer
	agent.Apis
	async.Task
}

func isMockCoreTaskRequest_Task2isTaskDetail_Task(in *gatewaypb.MockCoreTaskRequest) *corepb.TaskDetail {
	var retVal = corepb.TaskDetail{TaskId: in.TaskId}
	switch t := in.GetTask().(type) {
	case *gatewaypb.MockCoreTaskRequest_StartAgentRequest:
		retVal.Task = &corepb.TaskDetail_StartAgentRequest{StartAgentRequest: t.StartAgentRequest}
	case *gatewaypb.MockCoreTaskRequest_CreateAgentRequest:
		retVal.Task = &corepb.TaskDetail_CreateAgentRequest{CreateAgentRequest: t.CreateAgentRequest}
	case *gatewaypb.MockCoreTaskRequest_EditAgentRequest:
		retVal.Task = &corepb.TaskDetail_EditAgentRequest{EditAgentRequest: t.EditAgentRequest}
	case *gatewaypb.MockCoreTaskRequest_RemoveAgentRequest:
		retVal.Task = &corepb.TaskDetail_RemoveAgentRequest{RemoveAgentRequest: t.RemoveAgentRequest}
	case *gatewaypb.MockCoreTaskRequest_StopAgentRequest:
		retVal.Task = &corepb.TaskDetail_StopAgentRequest{StopAgentRequest: t.StopAgentRequest}
	case *gatewaypb.MockCoreTaskRequest_UpdateAgentRequest:
		retVal.Task = &corepb.TaskDetail_UpdateAgentRequest{UpdateAgentRequest: t.UpdateAgentRequest}
	case *gatewaypb.MockCoreTaskRequest_GetAgentStatusRequest:
		retVal.Task = &corepb.TaskDetail_GetAgentStatusRequest{GetAgentStatusRequest: t.GetAgentStatusRequest}
	}
	return &retVal
}
func (s gatewayDebugServer) MockCoreTask(_ context.Context, req *gatewaypb.MockCoreTaskRequest) (*gatewaypb.MockCoreTaskResponse, error) {
	if !req.GetIsQuery() {
		s.Task.AddTask(isMockCoreTaskRequest_Task2isTaskDetail_Task(req), true)
		return &gatewaypb.MockCoreTaskResponse{
			TaskId: req.TaskId,
			Result: &gatewaypb.MockCoreTaskResponse_Pending{},
		}, nil
	}

	t := s.Task.Query(req.TaskId)
	var retVal *gatewaypb.MockCoreTaskResponse
	switch result := t.GetResult().(type) {
	case *corepb.QueryTaskResultResponse_Finish:
		retVal = &gatewaypb.MockCoreTaskResponse{
			TaskId: req.TaskId,
			Result: &gatewaypb.MockCoreTaskResponse_Finish{Finish: result.Finish},
		}
	case *corepb.QueryTaskResultResponse_Failed:
		retVal = &gatewaypb.MockCoreTaskResponse{
			TaskId: req.TaskId,
			Result: &gatewaypb.MockCoreTaskResponse_Failed{Failed: result.Failed},
		}
	case *corepb.QueryTaskResultResponse_Working:
		retVal = &gatewaypb.MockCoreTaskResponse{
			TaskId: req.TaskId,
			Result: &gatewaypb.MockCoreTaskResponse_Working{Working: result.Working},
		}
	case *corepb.QueryTaskResultResponse_Pending:
		retVal = &gatewaypb.MockCoreTaskResponse{
			TaskId: req.TaskId,
			Result: &gatewaypb.MockCoreTaskResponse_Pending{Pending: result.Pending},
		}
	case *corepb.QueryTaskResultResponse_NotFound:
		retVal = &gatewaypb.MockCoreTaskResponse{
			TaskId: req.TaskId,
			Result: &gatewaypb.MockCoreTaskResponse_NotFound{NotFound: result.NotFound},
		}
	}
	return retVal, nil
}
func (s gatewayDebugServer) GatewayResources(_ context.Context,
	req *gatewaypb.GetGatewayResourcesApiRequest) (
	*gatewaypb.GetGatewayResourcesApiResponse, error) {
	switch r := req.GetResources().(type) {
	case *gatewaypb.GetGatewayResourcesApiRequest_Pool:
		res, _ := anypb.New(&gatewaypb.AgentPoolResources{AgentID: s.Apis.GetPoolIDs()})
		return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
	case *gatewaypb.GetGatewayResourcesApiRequest_AgentOperator:
		switch r := r.AgentOperator.Operator.(type) {
		//*AgentPoolOperator_StartAgentRequest
		//*AgentPoolOperator_CreateAgentRequest
		//todo:*AgentPoolOperator_EditAgentRequest
		//todo:*AgentPoolOperator_RemoveAgentRequest
		//*AgentPoolOperator_StopAgentRequest
		//todo:*AgentPoolOperator_UpdateAgentRequest
		//*AgentPoolOperator_GetAgentStatusRequest
		case *gatewaypb.AgentPoolOperator_StartAgentRequest:
			var allStatus status.Status
			for _, id := range r.StartAgentRequest.GetAgentId() {
				var subStatus status.Status
				err := s.Apis.StartAgent(id)
				if err != nil {
					subStatus.Code = 1
					subStatus.Message = err.Error()
				} else {
					subStatus.Message = "done"
				}
				result, _ := anypb.New(&subStatus)
				allStatus.Details = append(allStatus.Details, result)
			}
			res, _ := anypb.New(&allStatus)
			return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
		case *gatewaypb.AgentPoolOperator_CreateAgentRequest:
			req := model.ProxypbCreateAgentRequest2CreateAgentRequest(r.CreateAgentRequest)
			err := s.Apis.CreateAgent(req)
			if err != nil {
				res, _ := anypb.New(&status.Status{Code: 1, Message: err.Error()})
				return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
			}
			res, _ := anypb.New(&status.Status{Code: 0, Message: "done"})
			return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
		case *gatewaypb.AgentPoolOperator_StopAgentRequest:
			var allStatus status.Status
			for _, id := range r.StopAgentRequest.GetAgentId() {
				var subStatus status.Status
				err := s.Apis.StopAgent(id)
				if err != nil {
					subStatus.Code = 1
					subStatus.Message = err.Error()
				} else {
					subStatus.Message = "done"
				}
				result, _ := anypb.New(&subStatus)
				allStatus.Details = append(allStatus.Details, result)
			}
			res, _ := anypb.New(&allStatus)
			return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
		case *gatewaypb.AgentPoolOperator_GetAgentStatusRequest:
			var allStatus status.Status
			for _, id := range r.GetAgentStatusRequest.GetAgentId() {
				var subStatus status.Status
				subStatus.Message = s.Apis.GetAgentStatus(id).String()
				result, _ := anypb.New(&subStatus)
				allStatus.Details = append(allStatus.Details, result)
			}
			res, _ := anypb.New(&allStatus)
			return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
		default:
			return nil, grpcStatus.Errorf(codes.Unimplemented, "method GatewayResources not implemented")
		}
	case *gatewaypb.GetGatewayResourcesApiRequest_DataOperator:
		id := r.DataOperator.GetAgentID()
		pool := data.NewDataStoreApis(id)
		switch r := r.DataOperator.GetOperator().(type) {
		case *gatewaypb.DataPoolOperator_StoreDataRequest:
			err := pool.Store(r.StoreDataRequest)
			if err != nil {
				res, _ := anypb.New(&status.Status{Code: 1, Message: err.Error()})
				return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
			}
			res, _ := anypb.New(&status.Status{Message: "done"})
			return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
		case *gatewaypb.DataPoolOperator_GetDataRequest, *gatewaypb.DataPoolOperator_TopDataRequest:
			res, _ := anypb.New(pool.Top())
			return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
		case *gatewaypb.DataPoolOperator_PopDataRequest:
			res, _ := anypb.New(pool.Pop())
			return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
		case *gatewaypb.DataPoolOperator_CleanDataRequest:
			err := pool.Clean()
			if err != nil {
				res, _ := anypb.New(&status.Status{Code: 1, Message: err.Error()})
				return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
			}
			res, _ := anypb.New(&status.Status{Message: "done"})
			return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
		case *gatewaypb.DataPoolOperator_SizeDataRequest:
			res, _ := anypb.New(wrapperspb.Int32(int32(pool.Size())))
			return &gatewaypb.GetGatewayResourcesApiResponse{Resources: res}, nil
		}
	}
	return nil, grpcStatus.Errorf(codes.Unimplemented, "method GatewayResources not implemented")
}

func NewDebugServer(bind string) {
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.GetLoggerInterceptor(),
			middleware.GetRecovery(middleware.GetRegistry(middleware.GetSrvMetrics()))))
	time.Sleep(3 * time.Second)
	gatewaypb.RegisterDebugServiceServer(grpcServer, &gatewayDebugServer{Apis: agent.NewAgentApis(nil), Task: async.NewAsyncTask()})
	logger.Infoln("gateway debug-grpc-server start at: ", bind)
	panic(grpcServer.Serve(lis))
}
