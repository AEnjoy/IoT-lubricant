package data

import (
	"sync"

	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
)

type Apis interface {
	// Store and Push are the same API
	Store(*agent.DataMessage) error
	// Push and Store is the same API
	Push(*agent.DataMessage) error
	// Pop 从数据队列中取出数据并清理
	Pop() *core.Data
	// Top 从数据队列中取出数据不清理
	Top() *core.Data
	// Clean 清空数据队列
	Clean() error
}

var (
	agentDataStore sync.Map // string: agentID-> Apis
)

// NewDataStoreApis 初始化并(或)获取DataStoreApis 对象
func NewDataStoreApis(id string) Apis {
	actual, _ := agentDataStore.LoadOrStore(id, newData())
	return actual.(Apis)
}
func newData() Apis {
	return &data{
		data:       make([]*agent.DataMessage, 0),
		sendSignal: make(chan struct{}, 1),
	}
}
