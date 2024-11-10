package gateway

import (
	"context"
)

type agentCtrl struct {
	agentDevice <-chan []byte // /agentData/+agentID
	reg         <-chan []byte // Topic_AgentRegister
	ctx         context.Context
	ctrl        context.CancelFunc
}
