package _default

import (
	"fmt"
)

var (
	AgentDefaultBind = fmt.Sprintf(":%d", AgentGrpcPort)
)

const (
	AgentDefaultConfigFileName  = "agent.yaml"
	AgentDefaultOpenapiFileName = "api.json"
	AgentGrpcPort               = 5436
)
