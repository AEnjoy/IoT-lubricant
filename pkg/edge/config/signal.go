package config

import "context"

var (
	GatherSignal = make(chan context.Context)
	StopSignal   = make(chan context.Context)
	ConfigModify = make(chan struct{})
)
