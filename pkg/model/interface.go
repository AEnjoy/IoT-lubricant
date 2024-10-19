package model

import (
	"context"

	"gorm.io/gorm"
)

var _ GatewayDbOperator = (*GatewayDb)(nil)
var _ CoreDbOperator = (*CoreDb)(nil)

type CoreDbOperator interface {
	Begin() *gorm.DB
	Commit(txn *gorm.DB)
	Rollback(txn *gorm.DB)

	IsGatewayIdExists(string) bool
	StoreAgentGatherData(ctx context.Context, txn *gorm.DB, id, content string) error
	GetDataCleaner(id string) (*Clean, error)
	GetAgentInfo(id string) (*Agent, error)
}

type GatewayDbOperator interface {
	GetServerInfo() ServerInfo
	IsAgentIdExists(string) bool
	GetAllAgentId() []string
	RemoveAgent(...string) bool
	GetAgentReportCycle(string) int
	GetAgentGatherCycle(string) int
}
