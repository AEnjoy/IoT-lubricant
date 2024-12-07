package _default

import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
)

var (
	AgentDefaultBind = fmt.Sprintf(":%d", model.AgentGrpcPort)
)

const (
	AgentDefaultConfigFileName = "agent.yaml"
)
