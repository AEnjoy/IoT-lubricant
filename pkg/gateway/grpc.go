package gateway

import (
	"context"
	"io"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
	"google.golang.org/grpc"
)

var pong = &gateway.DataMessage{
	Flag: 0,
	Data: []byte("Pong"),
}

type grpcServer struct {
	gateway.UnimplementedGatewayServiceServer
}

func (*grpcServer) Data(stream grpc.BidiStreamingServer[gateway.DataMessage, gateway.DataMessage]) error {
	for {
		select {
		case <-stream.Context().Done():
			logger.Infof("Agent(ID:%s) disconnected", stream.Context().Value("agent_id"))
			return stream.Context().Err()
		case data := <-recvFromStream(stream):
			// 处理从客户端接收到的数据
			if data != nil {
				switch data.GetFlag() {
				case 2:
					dataRev <- (*model.EdgeData)(data)
					if err := stream.Send(pong); err != nil {
						return err
					}
				case 3:
					errMessages <- (*model.EdgeData)(data)
					// todo: 需要发送错误处理的响应结果
					fallthrough
				case 0: // Ping
					if err := stream.Send(pong); err != nil {
						return err
					}
				}
			}

		// 监听外部 ch1 通道是否有数据
		case msg := <-dataSend:
			// 处理并发送从 ch1 接收到的数据
			msg.Flag = 1
			if err := stream.Send((*gateway.DataMessage)(msg)); err != nil {
				return err
			}
		}
	}
}
func (*grpcServer) Ping(_ context.Context, ping *gateway.PingPong) (*gateway.PingPong, error) {
	if ping.GetFlag() == 0 {
		return &gateway.PingPong{Flag: 1}, nil
	}
	return &gateway.PingPong{Flag: 2}, nil
}
func (*grpcServer) PushMessageId(_ context.Context, in *gateway.AgentMessageIdInfo) (*gateway.AgentMessageIdInfo, error) {
	messageQueue <- in
	finish.Store(in.MessageId, make(chan struct{}))
	if v, ok := finish.Load(in.MessageId); ok {
		ch := v.(chan struct{})
		<-ch
		finish.Delete(in.MessageId)
		return &gateway.AgentMessageIdInfo{
			MessageId: in.MessageId,
			Time:      time.Now().Format("2006-01-02 15:04:05"),
		}, nil
	}
	return &gateway.AgentMessageIdInfo{}, nil
}

func recvFromStream(stream grpc.BidiStreamingServer[gateway.DataMessage, gateway.DataMessage]) <-chan *gateway.DataMessage {
	out := make(chan *gateway.DataMessage)

	go func() {
		defer close(out)
		for {
			data, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}

				// 返回 nil 以表示错误情况
				out <- nil
				return
			}
			out <- data
		}
	}()

	return out
}

func NewGrpcServer() *grpc.Server {
	middlewares := NewInterceptorImpl()
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middlewares.UnaryServerInterceptor),
		grpc.ChainStreamInterceptor(middlewares.StreamServerInterceptor),
	)
	gateway.RegisterGatewayServiceServer(server, &grpcServer{})
	return server
}
