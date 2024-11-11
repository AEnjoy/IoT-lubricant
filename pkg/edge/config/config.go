package config

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
)

var (
	Config    *types.EdgeSystem
	GatewayID string
)

func NullConfig() *types.EdgeSystem {
	return &types.EdgeSystem{
		Config:       &openapi.ApiInfo{},
		EnableConfig: &openapi.ApiInfo{},
	}
}
