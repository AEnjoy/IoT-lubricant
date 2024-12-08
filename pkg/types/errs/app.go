package errs

import "errors"

// Core
var (
	ErrTargetNoTask = errors.New("target has no task")
	ErrTimeout      = errors.New("get task timeout")
)

// Edge-Agent
var (
	ErrInvalidConfig      = errors.New("invalid config. Please check if all necessary settings have been set")
	ErrMultGatherInstance = errors.New("only one instance of the gather module can be started")
)

// Gateway-proxy
var (
	ErrAgentNotFound = errors.New("agent not found")
)
