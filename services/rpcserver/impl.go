package rpcserver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/types/errs"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	status2 "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (i PbCoreServiceImpl) Report(ctx context.Context, req *corepb.ReportRequest) (*corepb.ReportResponse, error) {
	gatewayid, _ := i.getGatewayID(ctx)
	userID, _ := i.getUserID(ctx)
	logg.L.Debugf("Recv gateway report request: %s", gatewayid)

	taskMq := i.DataStore.Mq
	data, err := proto.Marshal(req)
	if err != nil {
		logg.L.WithOperatorID(userID).Errorf("failed to marshal protobuf: %v", err)
		return nil, status2.Errorf(codes.Internal, "failed to marshal protobuf: %v", err)
	}

	err = taskMq.PublishBytes("/handler/report", data)
	if err != nil {
		logg.L.WithOperatorID(gatewayid).Errorf("failed to publish data to message queue: %v", err)
		return nil, status2.Errorf(codes.Internal, "failed to publish data: %v", err)
	}

	switch req.Req.(type) {
	case *corepb.ReportRequest_AgentStatus:
		return &corepb.ReportResponse{Resp: &corepb.ReportResponse_AgentStatus{
			AgentStatus: &corepb.AgentStatusResponse{
				Resp: &status.Status{Message: "ok"},
			},
		}}, nil
	case *corepb.ReportRequest_TaskResult:
		return &corepb.ReportResponse{Resp: &corepb.ReportResponse_TaskResult{
			TaskResult: &corepb.TaskResultResponse{
				Resp: &status.Status{Message: "ok"},
			},
		}}, nil
	case *corepb.ReportRequest_ReportLog:
		return &corepb.ReportResponse{Resp: &corepb.ReportResponse_ReportLog{
			ReportLog: &corepb.ReportLogResponse{
				Resp: &status.Status{Message: "ok"},
			},
		}}, nil
	case *corepb.ReportRequest_Error:
		return &corepb.ReportResponse{Resp: &corepb.ReportResponse_Error{
			Error: &corepb.ReportErrorResponse{
				Resp: &status.Status{Message: "ok"},
			},
		}}, nil
	default:
		return nil, status2.Errorf(codes.Internal, "unknown request type")
	}
}

func (i PbCoreServiceImpl) PushData(ctx context.Context, in *corepb.Data) (*corepb.PushDataResponse, error) {
	gatewayid, _ := i.getGatewayID(ctx)
	logger.Debugf("Recv data stream from gateway:%s", gatewayid)
	go i.handelRecvData(in)
	return &corepb.PushDataResponse{Resp: &status.Status{Code: 0, Message: "ok"}}, nil
}
func (i PbCoreServiceImpl) Ping(s grpc.BidiStreamingServer[metapb.Ping, metapb.Ping]) error {
	gatewayID, _ := i.getGatewayID(s.Context())
	userid, _ := i.getUserID(s.Context())
	taskMq := i.DataStore.Mq

	topic := fmt.Sprintf("/ping/%s/%s/register", taskTypes.TargetGateway, gatewayID)
	var (
		recSig  = make(chan struct{})
		exitSig = make(chan struct{})
		tryPing = false

		closeOnce sync.Once
	)
	defer func() {
		// clean
		closeOnce.Do(func() {
			close(exitSig)
			close(recSig)
		})

		// todo: if tryPing=true: send exception offline message to mq
		//  else send stand offline message to mq
		if tryPing {
			err := taskMq.Publish(fmt.Sprintf("/monitor/%s/%s/offline/error", taskTypes.TargetGateway, gatewayID), []byte(time.Now().Format("2006-01-02 15:04:05")))
			if err != nil {
				logg.L.WithOperatorID(userid).
					Errorf("failed to add gateway register information to messageQueue: %v gatewayID: %s", err, gatewayID)
			}
		} else {
			err := i.gatewayOffline(taskMq, userid, gatewayID)
			if err != nil {
				logg.L.WithOperatorID(userid).
					Errorf("failed to add gateway register information to messageQueue: %v gatewayID: %s", err, gatewayID)
			}
		}
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logg.L.WithOperatorID(userid).Errorf("recover error: %v", err)
			}
		}()
		for {
			resp, err := s.Recv()
			if err == io.EOF {
				closeOnce.Do(func() {
					close(exitSig)
					close(recSig)
				})
				return
			}
			if s.Context().Err() == context.Canceled {
				//exitSig <- struct{}{}
				return
			}

			if err != nil {
				logg.L.WithOperatorID(userid).Errorf("grpc stream error: %s", err.Error())
				continue
			}
			if tryPing && resp.Flag == 1 {
				tryPing = false
			}
			// logger.Debugf("Recv: Pong from Gateway: ID:%s", gatewayID)
			recSig <- struct{}{}
			time.Sleep(5 * time.Second)
		}
	}()
	for {
		select {
		case <-s.Context().Done():
			return nil
		case <-recSig:
			if err := s.Send(&metapb.Ping{Flag: 1}); err != nil {
				if s.Context().Err() == context.Canceled {
					//exitSig <- struct{}{}
					return nil
				}
				logg.L.WithOperatorID(userid).Errorf("grpc stream error: %s", err.Error())
				continue
			}
		case <-time.After(10 * time.Second):
			if !tryPing {
				tryPing = true
			} else {
				// time out
				return errors.New("gateway offline")
			}
		}

		err := taskMq.Publish(topic, []byte(gatewayID))
		if err != nil {
			logg.L.WithOperatorID(userid).Errorf("failed to add gateway register information to messageQueue: %v gatewayID: %s", err, gatewayID)
		}

	}
}
func (i PbCoreServiceImpl) GetTask(s grpc.BidiStreamingServer[corepb.Task, corepb.Task]) error {
	gatewayID, _ := i.getGatewayID(s.Context()) // 获取网关ID
	userid, _ := i.getUserID(s.Context())

	cancelContext, cancel := context.WithCancel(s.Context())
	defer cancel()
	taskMq := i.Mq

	ch, cancel2, err := i.getTaskIDCh(s.Context(), taskTypes.TargetGateway, userid, gatewayID)
	if err != nil {
		logg.L.WithOperatorID(userid).Errorf("failed to get task id: %s", err.Error())
		i.taskSendErrorMessage(s, 500, err.Error())
		return err
	}
	defer cancel2()

	// send core->gateway
	go func() {
		for {
			select {
			case <-cancelContext.Done():
				_ = taskMq.Publish(fmt.Sprintf("/monitor/%s/%s/offline", taskTypes.TargetGateway, gatewayID),
					[]byte(time.Now().Format("2006-01-02 15:04:05")))
				cancel2()
				return
			case idData := <-ch:
				taskID := idData
				if taskID == "" {
					continue
				}

				logg.L.Debugf("taskID:%s", taskID)

				taskData, err := i.getTask(s.Context(), taskTypes.TargetGateway, userid, gatewayID, taskID)
				if err != nil {
					logg.L.Debugf("Error at get task: %v", err)
					if err != errs.ErrTargetNoTask {
						logg.L.WithOperatorID(userid).Errorf("failed to get task id: %v", err)
						i.taskSendErrorMessage(s, 500, err.Error())
					}
				} else {
					logg.L.Debugf("send task %s to gateway %s", taskID, gatewayID)
					var resp corepb.CorePushTaskRequest
					var message corepb.TaskDetail
					err := proto.Unmarshal(taskData, &message)
					if err != nil {
						logg.L.WithOperatorID(userid).Errorf("failed to unmarshal task data: %v", err)
						continue
					}

					resp.Message = &message
					_ = s.Send(&corepb.Task{ID: taskID, Task: &corepb.Task_CorePushTaskRequest{CorePushTaskRequest: &resp}})
				}
			}
		}
	}()

	// recv gateway->core
	for {
		taskReq, err := s.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			// todo: handel error message more detail
			return err
		}
		// logger.Debugf("Recv: %v", taskReq)

		// HandelRecvGetTask(task) [Gateway->Core]
		switch taskReq.GetTask().(type) {
		case *corepb.Task_GatewayTryGetTaskRequest:
			// todo: 其实感觉没有必要，因为任务是由core主动获取并推送到gateway的，所以这里其实可以不用处理，直接返回即可
			//  clean needed
			targetID := gatewayID //task.GatewayTryGetTaskRequest.GetGatewayID()

			select {
			case taskID, ok := <-ch:
				if ok {
					logger.Debugf("taskID:%s", taskID)
					if taskData, err := i.getTask(s.Context(), taskTypes.TargetGateway, userid, targetID, taskID); err != nil {
						if err == errs.ErrTargetNoTask {
							//taskSendErrorMessage(s, 404, ErrTargetNoTask.Error())
							var resp corepb.NoTaskResponse
							var message corepb.TaskDetail
							message.TaskId = taskID
							//resp.Message = &message
							_ = s.Send(&corepb.Task{ID: taskReq.ID, Task: &corepb.Task_NoTaskResponse{NoTaskResponse: &resp}})
							continue
						} else {
							logg.L.WithOperatorID(userid).Errorf("failed to get task data: %s", err.Error())
							i.taskSendErrorMessage(s, 500, err.Error())
							continue
						}
					} else {
						var resp corepb.GatewayGetTaskResponse
						var message corepb.TaskDetail
						message.Content = taskData
						message.TaskId = taskID
						//resp.Message = &message
						logg.L.Debugf("send task to gateway: %s", taskID)
						_ = s.Send(&corepb.Task{ID: taskReq.ID, Task: &corepb.Task_GatewayGetTaskResponse{GatewayGetTaskResponse: &resp}})
						continue
					}
				}
			case <-time.After(1500 * time.Millisecond):
			}
			// 超时意味着没有创建过任务ch (任务不存在)
			// logger.Debugf("task not found")
			_ = s.Send(&corepb.Task{ID: taskReq.ID,
				Task: &corepb.Task_NoTaskResponse{NoTaskResponse: &corepb.NoTaskResponse{}},
			})
			//taskSendErrorMessage(s, 404, errs.ErrTimeout.Error())
		case *corepb.Task_CoreQueryTaskResultResponse:
			m := taskReq.GetCoreQueryTaskResultResponse()
			marshal, err := proto.Marshal(m)
			if err != nil {
				logg.L.WithOperatorID(userid).Errorf("failed to marshal task data: %v", err)
				continue
			}
			_ = taskMq.Publish(fmt.Sprintf("/task/%s/%s/%s/response", taskTypes.TargetGateway, gatewayID, m.TaskId), marshal)
		}
	}
}

// not impl
//
//	func (PbCoreServiceImpl) PushMessageId(context.Context, *corepb.MessageIdInfo) (*corepb.MessageIdInfo, error) {
//		return nil, nil
//	}
func (i PbCoreServiceImpl) PushDataStream(d grpc.BidiStreamingServer[corepb.Data, corepb.Data]) error {
	gatewayid, _ := i.getGatewayID(d.Context())
	userId, _ := i.getUserID(d.Context())
	logg.L.Debugf("Recv data stream from gateway:%s", gatewayid)

	for {
		data, err := d.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			// todo: handel error message more detail
			return err
		}
		// 由于数据处理需要消耗一定时间，所以使用goroutine处理
		go i.handelRecvData(data)

		err = d.Send(&corepb.Data{
			MessageId: data.MessageId,
		})
		if err != nil {
			logg.L.WithOperatorID(userId).Errorf("grpc stream error: %s", err.Error())
			continue
		}
	}
}
