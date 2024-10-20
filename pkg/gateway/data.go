package gateway

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
	"github.com/google/uuid"
)

const maxBuffer = 50

var (
	dataRev      = make(chan *types.EdgeData, maxBuffer)
	errMessages  = make(chan *types.EdgeData, maxBuffer)
	messageQueue = make(chan *gateway.AgentMessageIdInfo, maxBuffer)

	finish = sync.Map{}
)

type data struct {
	context.Context
	context.CancelFunc
}

func (d *data) Start() {
	d.Context, d.CancelFunc = context.WithCancel(context.Background())
	go d.HandleData(d.Context)
	go d.HandleError(d.Context)
}
func (d *data) Stop() {
	d.CancelFunc()
}
func (d *data) HandleError(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-errMessages:
			logger.Errorf("来自Agent(ID:%s)的错误:%s", e.AgentId, string(e.Data))
			// TODO: handle error
		}
	}
}
func (d *data) HandleData(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-dataRev:
			go d.handleData(data)
		// save data and send to core
		case data := <-messageQueue:
			go d.handleMessage(data)
		}
	}
}
func (d *data) handleData(in *types.EdgeData) {

}
func (d *data) handleMessage(in *gateway.AgentMessageIdInfo) {
	v, _ := finish.Load(in.MessageId)
	v.(chan struct{}) <- struct{}{} // 通知数据处理完成

	// TODO: handle message
}
func (a *app) handelSignal(id string) error {
	// todo: not all implemented yet
	//  task: 1. need to support choose report cycle
	//        2. need to support modify and trigger by core server manually
	v, ok := agentStore.Load(id)
	if !ok {
		return ErrAgentNotFound
	}
	agentMap := v.(*agentData)
	for range time.Tick(time.Second * 5) {
		agentMap.sendSignal <- struct{}{}
	}
	return nil
}
func (a *app) handelPushDataToServer(ctx context.Context, id string) error {
	v, ok := agentStore.Load(id)
	if !ok {
		return ErrAgentNotFound
	}
	agentMap := v.(*agentData)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-agentMap.sendSignal:
			go func() {
				_ = pushDataToServer(ctx, agentMap, a.grpcClient, id)
				// todo: this error need to be handled
			}()
		}
	}
}
func pushDataToServer(ctx context.Context, agentMap *agentData, grpcClient core.CoreServiceClient, id string) error {
	stream, err := grpcClient.PushData(ctx)
	if err != nil {
		return err
	}
	data := agentMap.coverToGrpcData()
	if data == nil {
		return nil
	}
	agentMap.cleanData()

	data.GatewayId = gatewayId
	data.AgentID = id
	data.MessageId = uuid.NewString()

	err = stream.Send(data)
	if err != nil {
		return err
	}

	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	if resp.GetMessageId() != data.MessageId {
		return errors.New("message id not match")
	}
	return nil
}
