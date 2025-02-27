package data

import (
	"sync"

	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
)

type Apis interface {
	// Store and Push are the same API
	Store(*agentpb.DataMessage) error
	// Push and Store is the same API
	Push(*agentpb.DataMessage) error
	// Pop 从数据队列中取出数据并清理
	Pop() *corepb.Data
	// Top 从数据队列中取出数据不清理
	Top() *corepb.Data
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
		data:       make([]*agentpb.DataMessage, 0),
		sendSignal: make(chan struct{}, 1),
	}
}
