package edge

import (
	"context"
	"io"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/exception"
	"github.com/AEnjoy/IoT-lubricant/pkg/gateway"
	pb "github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcClient(address, agentId string) (*grpc.ClientConn, error) {
	return grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(gateway.NewClientPerRPCCredentials(agentId)),
	)
}
func (a *app) clientGrpc() {
	stream, err := a.grpcClient.Data(context.Background())
	if err != nil {
		exception.ErrCh <- exception.New(1, "grpc clientGrpc error:", err.Error())
	}

	// recv
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				exception.ErrCh <- err
				return
			}
			switch in.Flag {
			// todo : need to handle
			}
		}
	}()

	// send
	for {
		select {
		case _ = <-dataSetCh:
			id := uuid.New()
			msg := &pb.DataMessage{
				AgentId: a.config.ID,
				//Data:      dataPacket, // todo: need to covert to []byte
				Flag:      2,
				MessageId: id.String(),
			}
			err := stream.Send(msg)
			if err != nil {
				exception.ErrCh <- err
			}

			go func() {
				resp, err := a.grpcClient.PushMessageId(context.Background(), &pb.AgentMessageIdInfo{
					MessageId: id.String(),
					AgentId:   a.config.ID,
					Time:      time.Now().Format("2006-01-02 15:04:05"),
				})
				if err != nil {
					exception.ErrCh <- err
				}
				if resp.MessageId != id.String() {
					exception.ErrCh <- exception.New(1, "grpc clientGrpc error:", "message id not match")
				}
			}()

		}
	}
}
