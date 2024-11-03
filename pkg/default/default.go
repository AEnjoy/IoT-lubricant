package _default

import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/pkg/types"
)

var (
	AgentDefaultBind = fmt.Sprintf(":%d", types.AgentGrpcPort)
)
