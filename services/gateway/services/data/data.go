package data

import (
	"math/rand"
	"sync"
	"sync/atomic"

	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
)

var _ Apis = (*data)(nil)

type data struct {
	data  []*agentpb.DataMessage
	cycle int //agent数据采集周期
	//loadTime   string // 第一次数据上载时间
	sendSignal chan struct{}
	l          sync.Mutex

	lastLen int
	cache   *corepb.Data
}

func (d *data) Size() int {
	return d.lastLen
}
func (d *data) Store(message *agentpb.DataMessage) error {
	return d.Push(message)
}

func (d *data) Push(message *agentpb.DataMessage) error {
	d.l.Lock()
	defer d.l.Unlock()

	d.parseData(message)
	return nil
}

func (d *data) Pop() *corepb.Data {
	d.l.Lock()
	defer d.l.Unlock()
	defer d.cleanData()

	return d.top()
}

func (d *data) Top() *corepb.Data {
	d.l.Lock()
	defer d.l.Unlock()

	return d.top()
}
func (d *data) top() *corepb.Data {
	if len(d.data) == 0 {
		return nil
	}

	if d.lastLen == len(d.data) && d.cache != nil {
		return d.cache
	} else {
		d.makeCache()
	}
	return d.cache
}

func (d *data) Clean() error {
	d.l.Lock()
	defer d.l.Unlock()

	d.cleanData()
	return nil
}

func (a *data) parseData(in *agentpb.DataMessage) {
	a.data = append(a.data, in)
	if rand.Intn(101) > 70 {
		a.makeCache()
	}
}
func (a *data) cleanData() {
	a.data = make([]*agentpb.DataMessage, 0)

	a.cache = nil
	a.lastLen = 0
}
func (a *data) coverToGrpcData() *corepb.Data {
	var data corepb.Data
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
func (d *data) makeCache() {
	d.cache = d.coverToGrpcData()
	d.lastLen = len(d.data)
}
