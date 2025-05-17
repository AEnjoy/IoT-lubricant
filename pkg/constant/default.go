package constant

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
const (
	USER_ID = "USER_ID"
)

// Whether it is a self hosted PaaS platform. If it needs to be provided to 'untrusted' tenants, please set it to true;
//
//	If the external tenant is trusted or used for internally, set it to false.
//	This will determine whether the tenant has the ability to upload scripts and run them.
const TrustedTenant = true
