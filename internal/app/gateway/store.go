package gateway

import (
	"sync"
	"sync/atomic"

	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
)

var agentDataStore sync.Map // string: agentID-> data: *agentData

type agentData struct {
	data  []*agent.DataMessage
	cycle int //agent数据采集周期
	//loadTime   string // 第一次数据上载时间
	sendSignal chan struct{}
	l          sync.Mutex
}

func (a *agentData) parseData(in *agent.DataMessage) {
	a.l.Lock()
	defer a.l.Unlock()
	a.data = append(a.data, in)
}
func (a *agentData) cleanData() {
	a.l.Lock()
	defer a.l.Unlock()
	a.data = make([]*agent.DataMessage, 0)
}
func (a *agentData) coverToGrpcData() *core.Data {
	a.l.Lock()
	defer a.l.Unlock()
	var data core.Data
	for _, datum := range a.data {
		for _, singleData := range datum.GetData() {
			data.Data = append(data.Data, singleData)
			atomic.AddInt32(&data.DataLen, 1)
		}
	}

	data.Cycle = int32(a.cycle)
	if len(a.data) > 0 {
		data.Time = a.data[0].GetDataGatherStartTime()
		return &data
	}
	return nil
}
