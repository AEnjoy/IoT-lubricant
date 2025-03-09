package response

type QueryMonitorBaseInfoResponse struct {
	GatewayCount int32 `json:"gatewayCount"`
	AgentCount   int32 `json:"agentCount"`
	NodeCount    int32 `json:"nodeCount"`

	OfflineGateway int32 `json:"offlineGateway"`
	OfflineAgent   int32 `json:"offlineAgent"`
}
