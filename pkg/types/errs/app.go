package errs

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
)

// Core
var (
	ErrTargetNoTask error = exception.New(code.ErrorCoreNoTask)
	ErrTimeout      error = exception.New(code.ErrorCoreTaskTimeout)
)

// Gateway-proxy
var (
	ErrAgentNotFound error = exception.New(code.ErrorGatewayAgentNotFound)
)

// Edge-Agent
var (
	ErrInvalidConfig      error = exception.New(code.ErrorAgentInvalidConfig, exception.WithMsg("Please check if all necessary settings have been set"))
	ErrMultGatherInstance error = exception.New(code.ErrorAgentNotAllowMultiGatherInstance, exception.WithMsg("Only one instance of the gather module can be started"))
)
