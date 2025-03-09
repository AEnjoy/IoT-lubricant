package lubricant

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/grpc/middleware"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/types"
	"github.com/aenjoy/iot-lubricant/pkg/types/errs"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	"github.com/aenjoy/iot-lubricant/services/lubricant/auth"
	"github.com/aenjoy/iot-lubricant/services/lubricant/config"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

var _ ioc.Object = (*Grpc)(nil)

// Grpc grpc server object client
type Grpc struct {
	*grpc.Server
	PbCoreServiceImpl
}

type PbCoreServiceImpl struct {
	corepb.UnimplementedCoreServiceServer
}

func (PbCoreServiceImpl) Report(ctx context.Context, req *corepb.ReportRequest) (*corepb.ReportResponse, error) {
	gatewayid, _ := getGatewayID(ctx)
	logger.Debugf("Recv gateway report request: %s", gatewayid)
	go HandelReport(req)

	if req.GetAgentStatus() != nil {
		return &corepb.ReportResponse{Resp: &corepb.ReportResponse_AgentStatus{
			AgentStatus: &corepb.AgentStatusResponse{
				Resp: &status.Status{Message: "ok"},
			},
		},
		}, nil
	}
	return &corepb.ReportResponse{Resp: &corepb.ReportResponse_Error{
		Error: &corepb.ReportErrorResponse{
			Resp: &status.Status{Message: "ok"},
		},
	},
	}, nil
}

func (PbCoreServiceImpl) PushData(ctx context.Context, in *corepb.Data) (*corepb.PushDataResponse, error) {
	gatewayid, _ := getGatewayID(ctx)
	logger.Debugf("Recv data stream from gateway:%s", gatewayid)
	go HandelRecvData(in)
	return &corepb.PushDataResponse{Resp: &status.Status{Code: 0, Message: "ok"}}, nil
}
func (PbCoreServiceImpl) Ping(s grpc.BidiStreamingServer[metapb.Ping, metapb.Ping]) error {
	gatewayID, _ := getGatewayID(s.Context())

	taskMq := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore).Mq

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
				logger.Errorf("failed to add gateway register information to messageQueue: %v gatewayID: %s", err, gatewayID)
			}
		} else {
			err := gatewayOffline(taskMq, gatewayID)
			if err != nil {
				logger.Errorf("failed to add gateway register information to messageQueue: %v gatewayID: %s", err, gatewayID)
			}
		}
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("recover error: %v", err)
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
				logger.Errorf("grpc stream error: %s", err.Error())
				continue
			}
			if tryPing && resp.Flag == 1 {
				tryPing = false
			}
			logger.Debugf("Recv: Pong from Gateway: ID:%s", gatewayID)
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
				logger.Errorf("grpc stream error: %s", err.Error())
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
			logger.Errorf("failed to add gateway register information to messageQueue: %v gatewayID: %s", err, gatewayID)
		}

	}
}
func (PbCoreServiceImpl) GetTask(s grpc.BidiStreamingServer[corepb.Task, corepb.Task]) error {
	gatewayID, _ := getGatewayID(s.Context()) // 获取网关ID
	cancelContext, cancel := context.WithCancel(s.Context())
	defer cancel()
	taskMq := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore).Mq

	ch, cancel2, err := getTaskIDCh(s.Context(), taskTypes.TargetGateway, gatewayID)
	if err != nil {
		logger.Errorf("failed to get task id: %s", err.Error())
		taskSendErrorMessage(s, 500, err.Error())
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

				logger.Debugf("taskID:%s", taskID)

				taskData, err := getTask(s.Context(), taskTypes.TargetGateway, gatewayID, taskID)
				if err != nil {
					logger.Debugf("Error at get task: %v", err)
					if err != errs.ErrTargetNoTask {
						logger.Errorf("failed to get task id: %v", err)
						taskSendErrorMessage(s, 500, err.Error())
					}
				} else {
					logger.Debugf("send task %s to gateway %s", taskID, gatewayID)
					var resp corepb.CorePushTaskRequest
					var message corepb.TaskDetail
					err := proto.Unmarshal(taskData, &message)
					if err != nil {
						logger.Errorf("failed to unmarshal task data: %v", err)
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
			_ = gatewayOffline(taskMq, gatewayID)
			return nil
		}
		if err != nil {
			// todo: handel error message more detail
			return err
		}
		logger.Debugf("Recv: %v", taskReq)

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
					if taskData, err := getTask(s.Context(), taskTypes.TargetGateway, targetID, taskID); err != nil {
						if err == errs.ErrTargetNoTask {
							//taskSendErrorMessage(s, 404, ErrTargetNoTask.Error())
							var resp corepb.NoTaskResponse
							var message corepb.TaskDetail
							message.TaskId = taskID
							//resp.Message = &message
							_ = s.Send(&corepb.Task{ID: taskReq.ID, Task: &corepb.Task_NoTaskResponse{NoTaskResponse: &resp}})
							continue
						} else {
							logger.Errorf("failed to get task data: %s", err.Error())
							taskSendErrorMessage(s, 500, err.Error())
							continue
						}
					} else {
						var resp corepb.GatewayGetTaskResponse
						var message corepb.TaskDetail
						message.Content = taskData
						message.TaskId = taskID
						//resp.Message = &message
						logger.Debugf("send task to gateway: %s", taskID)
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
				logger.Errorf("failed to marshal task data: %v", err)
				continue
			}
			_ = taskMq.Publish(fmt.Sprintf("/task/%s/%s/%s/response", taskTypes.TargetGateway, gatewayID, m.TaskId), marshal)
		}
	}
}
func (PbCoreServiceImpl) PushMessageId(context.Context, *corepb.MessageIdInfo) (*corepb.MessageIdInfo, error) {
	return nil, nil
}
func (PbCoreServiceImpl) PushDataStream(d grpc.BidiStreamingServer[corepb.Data, corepb.Data]) error {
	gatewayid, _ := getGatewayID(d.Context())
	logger.Debugf("Recv data stream from gateway:%s", gatewayid)
	mq := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore).Mq
	for {
		data, err := d.Recv()
		if err == io.EOF {
			_ = gatewayOffline(mq, gatewayid)
			return nil
		}
		if err != nil {
			// todo: handel error message more detail
			return err
		}
		// 由于数据处理需要消耗一定时间，所以使用goroutine处理
		go HandelRecvData(data)

		err = d.Send(&corepb.Data{
			MessageId: data.MessageId,
		})
		if err != nil {
			logger.Errorf("grpc stream error: %s", err.Error())
			continue
		}
	}
}

func (g *Grpc) Init() error {
	c := config.GetConfig()
	middlewares := ioc.Controller.Get(ioc.APP_NAME_CORE_GRPC_AUTH_INTERCEPTOR).(*auth.InterceptorImpl)
	var server *grpc.Server

	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     time.Minute * 10, // 连接空闲超过 10 分钟则关闭
		MaxConnectionAge:      time.Hour * 2,    // 连接存活超过 2 小时则关闭
		MaxConnectionAgeGrace: time.Minute * 5,  // 连接关闭前的宽限期
		Time:                  time.Minute * 1,  // 每隔 1 分钟发送一次 ping
		Timeout:               time.Second * 20, // ping 超时时间为 20 秒
	}
	kaep := keepalive.EnforcementPolicy{
		MinTime:             time.Minute * 5, // 连接建立后至少 5 分钟才能发送 ping
		PermitWithoutStream: true,            // 允许在没有流的情况下发送 ping
	}

	if c.Tls.Enable {
		grpcTlsOption, err := c.Tls.GetServerTlsConfig()
		if err != nil {
			return err
		}
		server = grpc.NewServer(
			grpcTlsOption,
			grpc.KeepaliveParams(kasp),
			grpc.KeepaliveEnforcementPolicy(kaep),
			grpc.MaxRecvMsgSize(1024*1024*100), // 100 MB
			grpc.MaxSendMsgSize(1024*1024*100), // 100 MB
			grpc.ChainStreamInterceptor(middlewares.StreamServerInterceptor),
			grpc.ChainUnaryInterceptor(
				//middlewares.UnaryServerInterceptor,
				middleware.GetLoggerInterceptor(),
				middleware.GetRecovery(middleware.GetRegistry(middleware.GetSrvMetrics())),
			),
		)
	} else {
		server = grpc.NewServer(
			grpc.KeepaliveParams(kasp),
			grpc.KeepaliveEnforcementPolicy(kaep),
			grpc.MaxRecvMsgSize(1024*1024*100), // 100 MB
			grpc.MaxSendMsgSize(1024*1024*100), // 100 MB
			grpc.ChainStreamInterceptor(middlewares.StreamServerInterceptor),
			grpc.ChainUnaryInterceptor(middlewares.UnaryServerInterceptor),
		)
	}
	corepb.RegisterCoreServiceServer(server, g)
	g.Server = server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.GrpcPort))
	if err != nil {
		return err
	}
	go func() {
		err := server.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
	logger.Debugf("core-grpc-server start at: %s", lis.Addr())
	return nil
}

func (Grpc) Weight() uint16 {
	return ioc.CoreGrpcServer
}

func (Grpc) Version() string {
	return "dev"
}

func taskSendErrorMessage(s grpc.BidiStreamingServer[corepb.Task, corepb.Task], code int, message string) {
	var out corepb.Task
	var errorResp corepb.Task_ErrorMessage
	errorResp.ErrorMessage = &metapb.ErrorMessage{Code: &status.Status{Code: int32(code), Message: message}}
	out.Task = &errorResp
	_ = s.Send(&out)
}
func gatewayOffline(mq mq.Mq, gatewayid string) error {
	logger.Debugf("gateway offline: %s", gatewayid)
	return mq.PublishBytes(fmt.Sprintf("/monitor/%s/offline", taskTypes.TargetGateway),
		[]byte(gatewayid))
}
func getGatewayID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata not found")
	}

	gatewayIDs := md.Get(string(types.NameGatewayID))
	if len(gatewayIDs) == 0 {
		return "", fmt.Errorf("gatewayid not found in metadata")
	}

	gatewayID := gatewayIDs[0]
	return gatewayID, nil
}
