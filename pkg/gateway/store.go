package gateway

import (
	"bytes"
	"sync"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/compress"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
)

var agentStore sync.Map // string: agentID-> data: *agentData

type agentData struct {
	data  []*gateway.DataMessage
	cycle int
	//loadTime   string // 第一次数据上载时间
	sendSignal chan struct{}
	l          sync.Mutex
}

func (a *agentData) parseData(in *gateway.DataMessage, cycle int) {
	a.l.Lock()
	defer a.l.Unlock()
	a.cycle = cycle
	t, _ := time.Parse("2006-01-02 15:04:05", in.Time)
	for i, data := range bytes.Split(in.Data, compress.Sepa) {
		a.data = append(a.data, &gateway.DataMessage{
			Data: data,
			Time: t.Add(time.Duration(i) * time.Second * time.Duration(cycle)).Format("2006-01-02 15:04:05"),
		})
	}
}
func (a *agentData) cleanData() {
	a.l.Lock()
	defer a.l.Unlock()
	a.data = make([]*gateway.DataMessage, 0)
}
func (a *agentData) coverToGrpcData() *core.Data {
	a.l.Lock()
	defer a.l.Unlock()
	var data core.Data
	data.DataLen = int32(len(data.Data))
	data.Cycle = int32(a.cycle)
	data.Time = a.data[0].GetTime()

	for _, datum := range a.data {
		data.Data = append(data.Data, datum.GetData())
	}

	return &data
}
