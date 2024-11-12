package config

import (
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
)

var (
	Config     *types.EdgeSystem
	GatewayID  string
	GatherLock sync.Mutex
)

func IsGathering() bool {
	if GatherLock.TryLock() {
		GatherLock.Unlock()
		return false
	}
	return true
}
func NullConfig() *types.EdgeSystem {
	return &types.EdgeSystem{
		Config:       &openapi.ApiInfo{},
		EnableConfig: &openapi.ApiInfo{},
	}
}
