package gateway

import (
	"fmt"
	"io"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	taskTypes "github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/AEnjoy/IoT-lubricant/protobuf/meta"
	json "github.com/bytedance/sonic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = &meta.Ping{
	Flag: 0,
}

func (a *app) grpcApp() error {
	// todo: not all implemented yet
	retryAttempts := 3        // 最大重试次数
	retryDelay := time.Second // 初始重试延迟

	for i := 0; i < retryAttempts; i++ {
		task, err := a.grpcClient.GetTask(a.ctrl)
		if err != nil {
			if i < retryAttempts-1 {
				time.Sleep(retryDelay)
				retryDelay *= 2 // 指数退避
				continue        // 重试
			}
			return fmt.Errorf("GetTask failed after %d attempts: %w", retryAttempts, err)
		}

		for {
			resp, err := task.Recv()
			if err == io.EOF {
				logger.Info("grpc stream closed")
				return nil
			}

			if err != nil {
				st, ok := status.FromError(err)
				if ok {
					logger.Errorf("grpc stream error: %s", st.Message())
					if st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded || st.Code() == codes.Aborted {
						logger.Errorf("Retrying after error (attempt %d): %v", i+1, err)
						if i < retryAttempts-1 {
							time.Sleep(retryDelay)
							retryDelay *= 2
							break
						}
						return fmt.Errorf("recv failed after %d attempts: %w", retryAttempts, err)
					} else {
						return fmt.Errorf("unrecoverable gRPC error: %w", err)
					}
				} else {
					logger.Errorf("Receive error: %v", err)
					return fmt.Errorf("recv failed: %w", err)
				}
				return err
			}

			var c types.TaskCommand
			var taskId string

			switch task := resp.GetTask().(type) {
			case *core.Task_GatewayGetTaskResponse:
				content := task.GatewayGetTaskResponse.GetMessage().GetContent()
				taskId = task.GatewayGetTaskResponse.GetMessage().GetTaskId()
				err = json.Unmarshal(content, &c)
				if err != nil {
					return err
				}
			case *core.Task_CorePushTaskRequest:
				content := task.CorePushTaskRequest.GetMessage().GetContent()
				taskId = task.CorePushTaskRequest.GetMessage().GetTaskId()
				err = json.Unmarshal(content, &c)
				if err != nil {
					return err
				}
			}

			switch c.ID {
			case taskTypes.OperationAddAgent:
				var req model.CreateAgentRequest
				err := json.Unmarshal(c.Data, &req)
				if err != nil {
					return err
				}

				response, err := HandelCreateAgentRequest(&req)
				if err != nil {
					return err
				}
				result, _ := json.Marshal(response)

				resp := &core.Task{
					ID: taskId,
					Task: &core.Task_CorePushTaskResponse{
						CorePushTaskResponse: &core.CorePushTaskResponse{
							Message: &core.TaskDetail{
								Content: result,
								TaskId:  taskId,
							},
						},
					},
				}
				_ = task.Send(resp)
			case taskTypes.OperationRemoveAgent:
				a.agentRemove("reserve a seat")
				panic("not implemented")
			case taskTypes.OperationNil:

			default:
				panic("unhandled default case")
			}
		}
	}
	return nil
}
