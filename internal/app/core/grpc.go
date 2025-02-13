package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/config"
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/internal/pkg/auth"
	"github.com/AEnjoy/IoT-lubricant/internal/pkg/grpc/middleware"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/errs"
	taskTypes "github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/AEnjoy/IoT-lubricant/protobuf/meta"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
)

var _ ioc.Object = (*Grpc)(nil)

// Grpc grpc server object client
type Grpc struct {
	*grpc.Server
	PbCoreServiceImpl
}

type PbCoreServiceImpl struct {
	core.UnimplementedCoreServiceServer
}

const timeOut = 1 // second
func createTimeOutContext(root context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(root, timeOut*time.Second)
}

func (PbCoreServiceImpl) Ping(s grpc.BidiStreamingServer[meta.Ping, meta.Ping]) error {
	gatewayID := s.Context().Value(types.NameGatewayID).(string)
	taskMq := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore).Mq

	topic := fmt.Sprintf("/ping/%s/%s/register", taskTypes.TargetGateway, gatewayID)
	var (
		recSig  = make(chan struct{})
		exitSig = make(chan struct{})
		tryPing = false
	)
	defer func() {
		// clean
		close(recSig)
		close(exitSig)
		// todo: if tryPing=true: send exception offline message to mq
		//  else send stand offline message to mq
	}()
	go func() {
		for {
			resp, err := s.Recv()
			if err == io.EOF {
				exitSig <- struct{}{}
				return
			}
			if err != nil {
				logger.Errorf("grpc stream error: %s", err.Error())
				continue
			}
			if tryPing && resp.Flag == 1 {
				tryPing = false
			}
			recSig <- struct{}{}
		}
	}()
	for {
		select {
		case <-s.Context().Done():
			return nil
		case <-recSig:
			if err := s.Send(&meta.Ping{Flag: 1}); err != nil {
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
func (PbCoreServiceImpl) GetTask(s grpc.BidiStreamingServer[core.Task, core.Task]) error {
	gatewayID := s.Context().Value(types.NameGatewayID).(string) // 获取网关ID
	cancelContext, cancel := context.WithCancel(s.Context())
	defer cancel()
	// send core->gateway
	taskMq := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore).Mq
	go func() {
		topic := fmt.Sprintf("/task/%s/%s", taskTypes.TargetGateway, gatewayID)
		taskIDCh, err := taskMq.Subscribe(topic)
		if err != nil {
			taskSendErrorMessage(s, 500, err.Error())
			return
		}
		for {
			select {
			case <-cancelContext.Done():
				_ = taskMq.Publish(fmt.Sprintf("/monitor/%s/%s/offline", taskTypes.TargetGateway, gatewayID),
					[]byte(time.Now().Format("2006-01-02 15:04:05")))
				if err := taskMq.Unsubscribe(topic, taskIDCh); err != nil {
					logger.Error("failed to unsubscribe task topic: %s", err.Error())
				}
				return
			case idData := <-taskIDCh:
				taskID := string(idData)
				if taskData, err := getTask(s.Context(), taskTypes.TargetGateway, gatewayID, taskID); err != nil {
					taskSendErrorMessage(s, 500, err.Error())
				} else {
					var resp core.CorePushTaskRequest
					var message core.TaskDetail
					message.Content = taskData
					message.TaskId = taskID
					resp.Message = &message
					_ = s.Send(&core.Task{ID: taskID, Task: &core.Task_CorePushTaskRequest{CorePushTaskRequest: &resp}})
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
			return err
		}

		// HandelRecvGetTask(task) [Gateway->Core]
		switch taskReq.GetTask().(type) {
		case *core.Task_GatewayTryGetTaskRequest:
			targetID := gatewayID //task.GatewayTryGetTaskRequest.GetGatewayID()
			ctx, cf := createTimeOutContext(s.Context())
			defer cf()

			ch, err := getTaskIDCh(ctx, taskTypes.TargetGateway, targetID)
			if err != nil {
				taskSendErrorMessage(s, 500, err.Error())
				continue
			}

			if taskID, ok := <-ch; ok {
				if taskData, err := getTask(ctx, taskTypes.TargetGateway, targetID, taskID); err != nil {
					if errors.Is(err, errs.ErrTargetNoTask) {
						//taskSendErrorMessage(s, 404, ErrTargetNoTask.Error())
						var resp core.NoTaskResponse
						var message core.TaskDetail
						message.TaskId = taskID
						//resp.Message = &message
						_ = s.Send(&core.Task{ID: taskReq.ID, Task: &core.Task_NoTaskResponse{NoTaskResponse: &resp}})
					} else {
						taskSendErrorMessage(s, 500, err.Error())
					}
				} else {
					var resp core.GatewayGetTaskResponse
					var message core.TaskDetail
					message.Content = taskData
					message.TaskId = taskID
					//resp.Message = &message
					_ = s.Send(&core.Task{ID: taskReq.ID, Task: &core.Task_GatewayGetTaskResponse{GatewayGetTaskResponse: &resp}})
				}
			} else {
				// 超时意味着没有创建过任务ch (任务不存在)
				taskSendErrorMessage(s, 500, errs.ErrTimeout.Error())
			}
		}
	}
}
func (PbCoreServiceImpl) PushMessageId(context.Context, *core.MessageIdInfo) (*core.MessageIdInfo, error) {
	return nil, nil
}
func (PbCoreServiceImpl) PushDataStream(d grpc.BidiStreamingServer[core.Data, core.Data]) error {
	for {
		data, err := d.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// 由于数据处理需要消耗一定时间，所以使用goroutine处理
		go HandelRecvData(data)

		err = d.Send(&core.Data{
			MessageId: data.MessageId,
		})
		if err != nil {
			return err
		}
	}
}

func (g *Grpc) Init() error {
	c := config.GetConfig()
	middlewares := ioc.Controller.Get(ioc.APP_NAME_CORE_GRPC_AUTH_INTERCEPTOR).(*auth.InterceptorImpl)
	var server *grpc.Server

	if c.Tls.Enable {
		grpcTlsOption, err := c.Tls.GetServerTlsConfig()
		if err != nil {
			return err
		}
		server = grpc.NewServer(
			grpcTlsOption,
			grpc.ChainStreamInterceptor(middlewares.StreamServerInterceptor),
			grpc.ChainUnaryInterceptor(
				middlewares.UnaryServerInterceptor,
				middleware.GetLoggerInterceptor(),
				middleware.GetRecovery(middleware.GetRegistry(middleware.GetSrvMetrics())),
			),
		)
	} else {
		server = grpc.NewServer(
			grpc.ChainStreamInterceptor(middlewares.StreamServerInterceptor),
			grpc.ChainUnaryInterceptor(middlewares.UnaryServerInterceptor),
		)
	}
	core.RegisterCoreServiceServer(server, g)
	g.Server = server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.GrpcPort))
	if err != nil {
		return err
	}
	go func() {
		_ = server.Serve(lis)
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

func taskSendErrorMessage(s grpc.BidiStreamingServer[core.Task, core.Task], code int, message string) {
	var out core.Task
	var errorResp core.Task_ErrorMessage
	errorResp.ErrorMessage = &meta.ErrorMessage{Code: &status.Status{Code: int32(code), Message: message}}
	out.Task = &errorResp
	_ = s.Send(&out)
}
