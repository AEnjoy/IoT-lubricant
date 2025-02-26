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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = &metapb.Ping{
	Flag: 0,
}

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

		// send - async task result
		go func() {
			for taskId := range a.task.GetNotifyCh() {
				result := a.task.Query(taskId).GetResult()
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
				}
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
						logger.Errorln("grpc stream closed")
						return
					}
					if err != nil {
						logger.Errorf("grpc stream error: %v", err)
					}
				}
			}
		}()
		// recv
		for {
			time.Sleep(time.Second)
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
			if resp.ID != "" {
				logger.Debugf("RecvTask %s:", resp.ID)
			}

			switch t := resp.GetTask().(type) {
			case *corepb.Task_GatewayGetTaskResponse:
				logger.Debug("Task Type is:", "GatewayGetTaskResponse")
				a.handelGatewayGetTaskResponse(t)
			case *corepb.Task_CorePushTaskRequest:
				logger.Debug("Task Type is:", "CorePushTaskRequest")
				resp := a.handelCorePushTaskAsync(t)
				err := task.Send(&corepb.Task{ID: resp.TaskId,
					Task: &corepb.Task_CorePushTaskResponse{
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
					logger.Debugf("Task %s finish", resp.GetTaskId())
					a.task.RemoveResult(resp.GetTaskId())
				}
			case *corepb.Task_GatewayQueryTaskResultResponse:
			// todo:
			case *corepb.Task_NoTaskResponse:
				logger.Infoln("gateway get task request success, and no task need to execute")
			case *corepb.Task_ErrorMessage:
				logger.Errorf("gateway send request to core success, but get an error: %s", t.ErrorMessage.String())
			}
		}
	}
	return nil
}
func (a *app) grpcDataApp() error {
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
					logger.Errorf("PushData exiting due to error: %v", err)
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
		a._checkPushDataStatus(resp)
	}
	return nil
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
			logger.Errorf("Failed to send ping request to server: %v", err)
			return err
		}
		for {
			if err := stream.Send(&metapb.Ping{Flag: 0}); err != nil {
				if err == io.EOF {
					logger.Errorln("grpc stream closed", "lost link with server")
					return nil
				}
				time.Sleep(time.Second)
				logger.Errorf("Failed to send ping request to server: %v", err)
				continue
			}
			_, err = stream.Recv()
			if err != nil {
				if err == io.EOF {
					logger.Errorln("grpc stream closed", "lost link with server")
					return nil
				}
				logger.Errorf("Failed to receive response from server: %v", err)
				return err
			}
		}

	}
	return nil
}
