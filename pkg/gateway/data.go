package gateway

import (
	"context"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
)

const maxBuffer = 50

var (
	dataSend     = make(chan *model.EdgeData, maxBuffer)
	dataRev      = make(chan *model.EdgeData, maxBuffer)
	errMessages  = make(chan *model.EdgeData, maxBuffer)
	messageQueue = make(chan *gateway.MessageIdInfo, maxBuffer)

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
	var dataModel uploadModel
	dataModel.edgeId = in.AgentId

}
func (d *data) handleMessage(in *gateway.MessageIdInfo) {
	v, _ := finish.Load(in.MessageId)
	v.(chan struct{}) <- struct{}{} // 通知数据处理完成

	// TODO: handle message
}
