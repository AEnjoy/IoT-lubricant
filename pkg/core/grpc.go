package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/auth"
	"github.com/AEnjoy/IoT-lubricant/pkg/core/config"
	"github.com/AEnjoy/IoT-lubricant/pkg/grpc/middleware"
	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	taskTypes "github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/AEnjoy/IoT-lubricant/protobuf/meta"
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

func (PbCoreServiceImpl) Ping(grpc.BidiStreamingServer[meta.Ping, meta.Ping]) error {
	return nil
}
func (PbCoreServiceImpl) GetTask(s grpc.BidiStreamingServer[core.Task, core.Task]) error {
	gatewayID := s.Context().Value(types.NameGatewayID).(string) // 获取网关ID
	// send core->gateway
	go func() {
		taskIDCh, err := taskMq.Subscribe(fmt.Sprintf("/task/%s/%s", taskTypes.TargetGateway, gatewayID))
		if err != nil {
			taskSendErrorMessage(s, 500, err.Error())
			return
		}
		for idData := range taskIDCh {
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
					if errors.Is(err, ErrTargetNoTask) {
						//taskSendErrorMessage(s, 404, ErrTargetNoTask.Error())
						var resp core.NoTaskResponse
						var message core.TaskDetail
						message.TaskId = taskID
						resp.Message = &message
						_ = s.Send(&core.Task{ID: taskReq.ID, Task: &core.Task_NoTaskResponse{NoTaskResponse: &resp}})
					} else {
						taskSendErrorMessage(s, 500, err.Error())
					}
				} else {
					var resp core.GatewayGetTaskResponse
					var message core.TaskDetail
					message.Content = taskData
					message.TaskId = taskID
					resp.Message = &message
					_ = s.Send(&core.Task{ID: taskReq.ID, Task: &core.Task_GatewayGetTaskResponse{GatewayGetTaskResponse: &resp}})
				}
			} else {
				// 超时意味着没有创建过任务ch (任务不存在)
				taskSendErrorMessage(s, 500, ErrTimeout.Error())
			}
		}
	}
}
func (PbCoreServiceImpl) PushMessageId(context.Context, *core.MessageIdInfo) (*core.MessageIdInfo, error) {
	return nil, nil
}
func (PbCoreServiceImpl) PushData(d grpc.BidiStreamingServer[core.Data, core.Data]) error {
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
	c := ioc.Controller.Get(config.APP_NAME).(*config.Config)
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
	errorResp.ErrorMessage = &meta.ErrorMessage{Code: int32(code), Message: message}
	out.Task = &errorResp
	_ = s.Send(&out)
}
