package main

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	def "github.com/aenjoy/iot-lubricant/pkg/default"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/types"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	"github.com/bytedance/sonic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func newCoreClient(address, userID, gatewayID string) (corepb.CoreServiceClient, context.Context, error) {
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*100), grpc.MaxCallSendMsgSize(1024*1024*100)), // 100 MB
	)
	if err != nil {
		return nil, context.TODO(), err
	}

	cli := corepb.NewCoreServiceClient(conn)
	md := metadata.New(map[string]string{
		string(types.NameGatewayID): gatewayID,
		def.USER_ID:                 userID,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := cli.Ping(ctx)
	if err != nil {
		logger.Errorf("Failed to send ping request to server: %v", err)
		return nil, context.TODO(), err
	}

	if err := stream.Send(&metapb.Ping{Flag: 0}); err != nil {
		logger.Errorf("Failed to send ping request to server: %v", err)
		return nil, context.TODO(), err
	}

	resp, err := stream.Recv()
	if err != nil {
		logger.Errorf("Failed to receive response from server: %v", err)
		return nil, context.TODO(), err
	}
	if resp.GetFlag() != 1 {
		return nil, context.TODO(), errors.New("lubricant server not ready")
	}
	return cli, ctx, stream.CloseSend()
}

var sendCountSuccess int32
var sendCountFail int32

func pushData2Core(cli corepb.CoreServiceClient, ctx context.Context, dataCh chan *Data) {
	var sendBufferSig = make(chan [][]byte)
	go func() {
		<-startSig

		for d := range sendBufferSig {
			_,err:=cli.PushData(ctx, &corepb.Data{
				GatewayId: gatewayID,
				AgentID:   randGetAgentID(),
				Data:      d,
				DataLen:   10,
				Time:time.Now().Add(-10*time.Second).Format("2006-01-02 15:04:05"),
				Cycle: 1,
			})
			if err!=nil{
				atomic.AddInt32(&sendCountFail, 1)
			}else{
				atomic.AddInt32(&sendCountSuccess, 1)
			}
		}
	}()

	var sendBuffer [][]byte
	var sendBufferCount int
	for data := range dataCh {
		dataBytes, err := sonic.Marshal(data)
		if err != nil {
			continue
		}

		sendBuffer = append(sendBuffer, dataBytes)
		sendBufferCount++

		if sendBufferCount == 10 {
			sendBufferSig <- sendBuffer
			sendBuffer = nil
			sendBufferCount = 0
		}
	}
}
