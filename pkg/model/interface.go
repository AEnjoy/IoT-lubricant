package model

var _ GatewayDbCli = (*GatewayDb)(nil)

type GatewayDbCli interface {
	GetServerInfo() ServerInfo
	IsAgentIdExists(string) bool
	GetAllAgentId() []string
	RemoveAgent(...string) bool
	GetAgentReportCycle(string) int
	GetAgentGatherCycle(string) int
}
