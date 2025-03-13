package response

type PushAgentTaskResponse struct {
	TaskID string `json:"task_id"`
}
type AddAgentResponse struct {
	AgentID string `json:"agent_id"`
	PushAgentTaskResponse
}
type AgentAsyncExecuteOperatorResponse struct {
	TaskID string `json:"taskId"`
	Data   string `json:"data,omitempty"`
}
type GetOpenApiDocResponse struct {
	AgentID string `json:"agentId"`

	Doc []byte `json:"doc"` // base64 encode it is openapi.ApiInfo
}
