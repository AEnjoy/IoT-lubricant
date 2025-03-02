package response

type QueryMonitorBaseInfoResponse struct {
	GatewayCount int `json:"gatewayCount"`
	AgentCount   int `json:"agentCount"`
	NodeCount    int `json:"nodeCount"`

	OfflineGateway int `json:"offlineGateway"`
	OfflineAgent   int `json:"offlineAgent"`
}
