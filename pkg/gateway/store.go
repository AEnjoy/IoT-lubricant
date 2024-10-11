package gateway

import (
	"sync"

	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
)

var agentStore sync.Map // string: agentID-> data: *agent

type agent struct {
	data       []gateway.DataMessage
	loadTime   string // 第一次数据上载时间
	sendSignal chan struct{}
	l          sync.Mutex
}

func (a *agent) parseData(in gateway.DataMessage) {

}
