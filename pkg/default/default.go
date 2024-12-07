package _default

import (
	"fmt"
)

var (
	AgentDefaultBind = fmt.Sprintf(":%d", AgentGrpcPort)
)

const (
	AgentDefaultConfigFileName = "agent.yaml"
	AgentGrpcPort              = 5436
)
