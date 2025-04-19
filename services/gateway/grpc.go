package gateway

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	"github.com/aenjoy/iot-lubricant/services/gateway/services/data"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _reportMessage = make(chan *corepb.ReportRequest, 100)

func (a *app) grpcReportApp() {
	for request := range _reportMessage {
		request.GatewayId = gatewayId
		_, err := a.grpcClient.Report(a.ctrl, request)
		if err != nil {
			logger.Errorf("Failed to send report request to server: %v", err)
		}
	}
}

// todo: 需要重构逻辑，以后Core向网关发送任务后，不需要再通过这个接口查询结果，而是通过网关主动Report上传结果
func (a *app) grpcTaskApp() error {
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
		i--

		// send - async task result
		go func() {
			for taskId := range a.task.GetNotifyCh() {
				logger.Debug("Task:", taskId)
				/*result := a.task.Query(taskId).GetResult()
				err := task.Send(&corepb.Task{
					ID: taskId,
					Task: &corepb.Task_CoreQueryTaskResultResponse{
						CoreQueryTaskResultResponse: &corepb.QueryTaskResultResponse{
							TaskId: taskId,
							Result: result,
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
				if _, ok := result.(*corepb.QueryTaskResultResponse_Finish); ok {
					logger.Debugf("Task %s finish", taskId)
					a.task.RemoveResult(taskId)
				}*/
			}
		}()

		// send -
		go func() {
			for {
				select {
				case <-a.ctrl.Done():
					return
				case <-time.After(5 * time.Second):
					err := task.Send(&corepb.Task{
						Task: &corepb.Task_GatewayTryGetTaskRequest{
							GatewayTryGetTaskRequest: &corepb.GatewayTryGetTaskRequest{
								GatewayID: gatewayId,
							},
						},
					})
					if err == io.EOF {
						logg.L.Error("grpc stream closed")
						return
					}
					if err != nil {
						logg.L.Errorf("grpc stream error: %v", err)
					}
				}
			}
		}()
		// recv
		for {
			time.Sleep(time.Second)
			resp, err := task.Recv()
			if err == io.EOF {
				logg.L.Error("grpc stream closed")
				return nil
			}

			if err != nil {
				st, ok := status.FromError(err)
				if ok {
					logg.L.Errorf("grpc stream error: %s", st.Message())
					if st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded || st.Code() == codes.Aborted {
						logg.L.Errorf("Retrying after error (attempt %d): %v", i+1, err)
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
					logg.L.Errorf("Receive error: %v", err)
					return fmt.Errorf("recv failed: %w", err)
				}
			}
			if resp.ID != "" {
				logg.L.Debugf("RecvTask %s:", resp.ID)
			}

			switch t := resp.GetTask().(type) {
			case *corepb.Task_GatewayGetTaskResponse:
				logg.L.Debug("Task Type is:", "GatewayGetTaskResponse")
				a.handelGatewayGetTaskResponse(t)
			case *corepb.Task_CorePushTaskRequest:
				logg.L.Debug("Task Type is:", "CorePushTaskRequest")
				_ = a.handelCorePushTaskAsync(t)
				//err := task.Send(&corepb.Task{ID: resp.TaskId,
				//	Task: &corepb.Task_CorePushTaskResponse{
				//		CorePushTaskResponse: resp,
				//	},
				//})
				//if err != nil {
				//	st, ok := status.FromError(err)
				//	if ok {
				//		if st.Code() == codes.Unavailable || st.Code() == codes.Canceled || st.Code() == codes.DeadlineExceeded || err == io.EOF {
				//			logger.Errorf("Send goroutine exiting due to error: %v", err)
				//			continue
				//		} else {
				//			logger.Errorf("Unrecoverable gRPC Send error: %v", err)
				//		}
				//	} else {
				//		logger.Errorf("Send error: %v", err)
				//		if errors.Is(err, io.EOF) {
				//			logger.Errorln("grpc stream closed")
				//			continue
				//		}
				//		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				//			logger.Errorln("net timeout")
				//		}
				//		continue
				//	}
				//	continue
				//}
			case *corepb.Task_CoreQueryTaskResultRequest:
				resp := a.task.Query(t.CoreQueryTaskResultRequest.GetTaskId())
				err := task.Send(&corepb.Task{ID: resp.TaskId,
					Task: &corepb.Task_CoreQueryTaskResultResponse{
						CoreQueryTaskResultResponse: resp,
					},
				})
				if err != nil {
					st, ok := status.FromError(err)
					if ok {
						if st.Code() == codes.Unavailable || st.Code() == codes.Canceled || st.Code() == codes.DeadlineExceeded || err == io.EOF {
							logg.L.Errorf("Send goroutine exiting due to error: %v", err)
							continue
						} else {
							logg.L.Errorf("Unrecoverable gRPC Send error: %v", err)
						}
					} else {
						logg.L.Errorf("Send error: %v", err)
						if errors.Is(err, io.EOF) {
							logg.L.Error("grpc stream closed")
							continue
						}
						if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
							logg.L.Error("net timeout")
						}
						continue
					}
					continue
				}
				//if resp.GetFinish() != nil {
				//	logger.Debugf("Task %s finish", resp.GetTaskId())
				//	a.task.RemoveResult(resp.GetTaskId())
				//}
			case *corepb.Task_GatewayQueryTaskResultResponse:
			// todo:
			case *corepb.Task_NoTaskResponse:
				// logger.Debug("gateway get task request success, and no task need to execute")
			case *corepb.Task_ErrorMessage:
				logg.L.Error("gateway send request to core success, but get an error: %s", t.ErrorMessage.String())
			}
		}
	}
	return nil
}
func (a *app) grpcDataApp() {
	time.Sleep(30 * time.Second)
	data.InitDataSendQueue()
	ch := data.GetDataSendQueue()
	for d := range ch {
		//todo:可选: 发送前需要 GetCoreCapacity 检查
		d.GatewayId = gatewayId
		resp, err := a.grpcClient.PushData(a.ctrl, d)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				if st.Code() == codes.Unavailable || st.Code() == codes.Canceled || st.Code() == codes.DeadlineExceeded || err == io.EOF {
					logg.L.Errorf("PushData exiting due to error: %v", err)
					continue
				} else {
					logg.L.Errorf("Unrecoverable gRPC Send error: %v", err)
				}
			} else {
				logg.L.Errorf("Send error: %v", err)
				if errors.Is(err, io.EOF) {
					logg.L.Error("grpc stream closed")
					continue
				}
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					logg.L.Error("net timeout")
				}
				continue
			}
			continue
		}
		a._checkPushDataStatus(resp)
	}
}
func (a *app) _checkPushDataStatus(resp *corepb.PushDataResponse) {
	// todo:
}
func (a *app) grpcPingApp() error {
	retryAttempts := 3        // 最大重试次数
	retryDelay := time.Second // 初始重试延迟
	for i := 0; i < retryAttempts; i++ {
		stream, err := a.grpcClient.Ping(a.ctrl)
		if err != nil {
			if i < retryAttempts-1 {
				time.Sleep(retryDelay)
				retryDelay *= 2 // 指数退避
				continue        // 重试
			}
			logg.L.Errorf("Failed to send ping request to server: %v", err)
			return err
		}
		i--

		for {
			if err := stream.Send(&metapb.Ping{Flag: 0}); err != nil {
				if err == io.EOF {
					logg.L.Error("grpc stream closed", "lost link with server")
					break
				}
				time.Sleep(time.Second)
				logg.L.Error("Failed to send ping request to server: %v", err)
				continue
			}
			_, err = stream.Recv()
			if err != nil {
				if err == io.EOF {
					logg.L.Error("grpc stream closed", "lost link with server")
					break
				}
				logg.L.Errorf("Failed to receive response from server: %v", err)
			}
		}
	}
	return nil
}
