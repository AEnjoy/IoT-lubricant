package errs

import (
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
)

// Core
var (
	ErrTargetNoTask error = exception.New(exceptionCode.ErrorCoreNoTask)
	ErrTimeout      error = exception.New(exceptionCode.ErrorCoreTaskTimeout)
)

// Gateway-gateway
var (
	ErrAgentNotFound error = exception.New(exceptionCode.ErrorGatewayAgentNotFound)
)

// Edge-Agent
var (
	ErrInvalidConfig      error = exception.New(exceptionCode.ErrorAgentInvalidConfig, exception.WithMsg("Please check if all necessary settings have been set"))
	ErrMultGatherInstance error = exception.New(exceptionCode.ErrorAgentNotAllowMultiGatherInstance, exception.WithMsg("Only one instance of the gather module can be started"))
)
