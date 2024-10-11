package gateway

import (
	"context"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
	"github.com/google/uuid"
)

const maxBuffer = 50

var (
	dataRev      = make(chan *model.EdgeData, maxBuffer)
	errMessages  = make(chan *model.EdgeData, maxBuffer)
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
func (d *data) handleData(in *model.EdgeData) {

}
func (d *data) handleMessage(in *gateway.AgentMessageIdInfo) {
	v, _ := finish.Load(in.MessageId)
	v.(chan struct{}) <- struct{}{} // 通知数据处理完成

	// TODO: handle message
}
func (a *app) pushDataToServer(ctx context.Context, id string) error {
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
			stream, err := a.grpcClient.PushData(ctx)
			if err != nil {
				return err
			}
			data := agentMap.coverToGrpcData()
			agentMap.cleanData()

			data.GatewayId = gatewayId
			data.AgentID = id
			data.MessageId = uuid.NewString()

			err = stream.Send(data)
			if err != nil {
				return err
			}
		}
	}
}
