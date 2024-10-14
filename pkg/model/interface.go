package model

var _ GatewayDbCli = (*GatewayDb)(nil)
var _ CoreDbCli = (*CoreDb)(nil)

type CoreDbCli interface {
	IsGatewayIdExists(string) bool
}

type GatewayDbCli interface {
	GetServerInfo() ServerInfo
	IsAgentIdExists(string) bool
	GetAllAgentId() []string
	RemoveAgent(...string) bool
	GetAgentReportCycle(string) int
	GetAgentGatherCycle(string) int
}
