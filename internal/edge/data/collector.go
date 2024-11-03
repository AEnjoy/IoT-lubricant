package data

import (
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
)

var DataCollect = make([]*edge.DataPacket, 0)

var DCL = sync.Mutex{}
