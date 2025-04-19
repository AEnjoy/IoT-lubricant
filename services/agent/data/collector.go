package data

import (
	"sync"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/edge"
)

var Collector Collect = newCollectSlot()

type Collect interface {
	AddData(slot int, data *edge.DataPacket)
	GetData(slot int) []*edge.DataPacket
	GetDataLen(slot int) int
	AddSlot(slot int)
	CleanData(slot int)
	RemoveSlot(slot int)
}

var _ Collect = (*collectSlot)(nil)

type collectSlot struct {
	dataCollect map[int][]*edge.DataPacket
	lock        map[int]*sync.Mutex // write dataCollect lock

	l sync.RWMutex
}

func (c *collectSlot) AddData(slot int, data *edge.DataPacket) {
	// todo:实际上，这里应该在第一次启动时便初始化
	if _, ok := c.dataCollect[slot]; !ok {
		c.AddSlot(slot)
	}

	c.lock[slot].Lock()
	defer c.lock[slot].Unlock()
	if len(c.dataCollect[slot]) == 0 {
		data.Timestamp = time.Now()
	}
	c.dataCollect[slot] = append(c.dataCollect[slot], data)
}

func (c *collectSlot) GetData(slot int) []*edge.DataPacket {
	c.l.RLock()
	d := c.dataCollect[slot]
	c.l.RUnlock()
	c.CleanData(slot)
	return d
}
func (c *collectSlot) CleanData(slot int) {
	c.l.Lock()
	defer c.l.Unlock()
	c.dataCollect[slot] = make([]*edge.DataPacket, 0)
}
func (c *collectSlot) GetDataLen(slot int) int {
	c.l.RLock()
	defer c.l.RUnlock()
	v, ok := c.dataCollect[slot]
	if !ok {
		return -1
	}
	return len(v)
}

func (c *collectSlot) AddSlot(slot int) {
	c.l.Lock()
	defer c.l.Unlock()
	c.dataCollect[slot] = make([]*edge.DataPacket, 0)
	c.lock[slot] = &sync.Mutex{}
}

func (c *collectSlot) RemoveSlot(slot int) {
	c.l.Lock()
	defer c.l.Unlock()
	delete(c.dataCollect, slot)
	delete(c.lock, slot)
}
func newCollectSlot() *collectSlot {
	return &collectSlot{
		dataCollect: make(map[int][]*edge.DataPacket),
		lock:        make(map[int]*sync.Mutex),
	}
}
