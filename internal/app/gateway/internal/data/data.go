package data

import (
	"math/rand"
	"sync"
	"sync/atomic"

	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
)

var _ Apis = (*data)(nil)

type data struct {
	data  []*agent.DataMessage
	cycle int //agent数据采集周期
	//loadTime   string // 第一次数据上载时间
	sendSignal chan struct{}
	l          sync.Mutex

	lastLen int
	cache   *core.Data
}

func (d *data) Store(message *agent.DataMessage) error {
	return d.Push(message)
}

func (d *data) Push(message *agent.DataMessage) error {
	d.l.Lock()
	defer d.l.Unlock()

	d.parseData(message)
	return nil
}

func (d *data) Pop() *core.Data {
	d.l.Lock()
	defer d.l.Unlock()
	defer d.cleanData()

	return d.top()
}

func (d *data) Top() *core.Data {
	d.l.Lock()
	defer d.l.Unlock()

	return d.top()
}
func (d *data) top() *core.Data {
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

func (a *data) parseData(in *agent.DataMessage) {
	a.data = append(a.data, in)
	if rand.Intn(101) > 70 {
		a.makeCache()
	}
}
func (a *data) cleanData() {
	a.data = make([]*agent.DataMessage, 0)

	a.cache = nil
	a.lastLen = 0
}
func (a *data) coverToGrpcData() *core.Data {
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
func (d *data) makeCache() {
	d.cache = d.coverToGrpcData()
	d.lastLen = len(d.data)
}
