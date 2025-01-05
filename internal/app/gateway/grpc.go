package gateway

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/AEnjoy/IoT-lubricant/protobuf/meta"
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

	// task
	for i := 0; i < retryAttempts; i++ {
		task, err := a.grpcClient.GetTask(a.ctrl)
		if err != nil {
			if i < retryAttempts-1 {
				time.Sleep(retryDelay)
				retryDelay *= 2 // 指数退避
				continue        // 重试
			}
			return fmt.Errorf("GetTask failed after %d attempts: %w\n", retryAttempts, err)
		}

		// send - async task result
		go func() {
			for taskId := range a.task.GetNotifyCh() {
				err := task.Send(&core.Task{
					ID: taskId,
					Task: &core.Task_CoreQueryTaskResultResponse{
						CoreQueryTaskResultResponse: &core.QueryTaskResultResponse{
							TaskId: taskId,
							Result: a.task.Query(taskId).GetResult(),
						},
					},
				})
				if err != nil {
					st, ok := status.FromError(err)
					if ok {
						if st.Code() == codes.Unavailable || st.Code() == codes.Canceled || st.Code() == codes.DeadlineExceeded || err == io.EOF {
							logger.Errorf("Send goroutine exiting due to error: %v", err)
							return
						} else {
							logger.Errorf("Unrecoverable gRPC Send error: %v", err)
						}
					} else {
						logger.Errorf("Send error: %v", err)
						if errors.Is(err, io.EOF) {
							logger.Errorln("grpc stream closed")
							return
						}
						if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
							logger.Errorln("net timeout")
						}
						return
					}
					return
				}
				a.task.RemoveResult(taskId)
			}
		}()

		// recv{

		for {
			resp, err := task.Recv()
			if err == io.EOF {
				logger.Errorln("grpc stream closed")
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
			}

			switch t := resp.GetTask().(type) {
			case *core.Task_GatewayGetTaskResponse:
				a.handelGatewayGetTaskResponse(t)
			case *core.Task_CorePushTaskRequest:
				resp := a.handelCorePushTaskAsync(t)
				err := task.Send(&core.Task{ID: resp.TaskId,
					Task: &core.Task_CorePushTaskResponse{
						CorePushTaskResponse: resp,
					},
				})
				if err != nil {
					st, ok := status.FromError(err)
					if ok {
						if st.Code() == codes.Unavailable || st.Code() == codes.Canceled || st.Code() == codes.DeadlineExceeded || err == io.EOF {
							logger.Errorf("Send goroutine exiting due to error: %v", err)
							continue
						} else {
							logger.Errorf("Unrecoverable gRPC Send error: %v", err)
						}
					} else {
						logger.Errorf("Send error: %v", err)
						if errors.Is(err, io.EOF) {
							logger.Errorln("grpc stream closed")
							continue
						}
						if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
							logger.Errorln("net timeout")
						}
						continue
					}
					continue
				}
				a.task.RemoveResult(resp.GetTaskId())
			case *core.Task_CoreQueryTaskResultRequest:
				resp := a.task.Query(t.CoreQueryTaskResultRequest.GetTaskId())
				err := task.Send(&core.Task{ID: resp.TaskId,
					Task: &core.Task_CoreQueryTaskResultResponse{
						CoreQueryTaskResultResponse: resp,
					},
				})
				if err != nil {
					st, ok := status.FromError(err)
					if ok {
						if st.Code() == codes.Unavailable || st.Code() == codes.Canceled || st.Code() == codes.DeadlineExceeded || err == io.EOF {
							logger.Errorf("Send goroutine exiting due to error: %v", err)
							continue
						} else {
							logger.Errorf("Unrecoverable gRPC Send error: %v", err)
						}
					} else {
						logger.Errorf("Send error: %v", err)
						if errors.Is(err, io.EOF) {
							logger.Errorln("grpc stream closed")
							continue
						}
						if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
							logger.Errorln("net timeout")
						}
						continue
					}
					continue
				}
				if resp.GetFinish() != nil {
					a.task.RemoveResult(resp.GetTaskId())
				}
			case *core.Task_GatewayQueryTaskResultResponse:
			// todo:
			case *core.Task_NoTaskResponse:
				logger.Infoln("gateway get task request success, and no task need to execute")
			case *core.Task_ErrorMessage:
				logger.Errorf("gateway send request to core success,but get the error: %s", t.ErrorMessage.String())
			}

			//switch c.ID {
			//case taskTypes.OperationAddAgent:
			//	var req model.CreateAgentRequest
			//	err := json.Unmarshal(c.Data, &req)
			//	if err != nil {
			//		return err
			//	}
			//
			//	response, err := HandelCreateAgentRequest(&req)
			//	if err != nil {
			//		return err
			//	}
			//	_, _ = json.Marshal(response)
			//
			//	resp := &core.Task{
			//		ID: taskId,
			//		Task: &core.Task_CorePushTaskResponse{
			//			CorePushTaskResponse: &core.CorePushTaskResponse{
			//				//Message: &core.TaskDetail{
			//				//	Content: result,
			//				//	TaskId:  taskId,
			//				//},
			//			},
			//		},
			//	}
			//	_ = task.Send(resp)
			//case taskTypes.OperationRemoveAgent:
			//	a.agentRemove("reserve a seat")
			//	panic("not implemented")
			//case taskTypes.OperationNil:
			//
			//default:
			//	panic("unhandled default case")
			//}
		}
	}
	return nil
}
