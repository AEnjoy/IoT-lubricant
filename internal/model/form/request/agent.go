package request

type AddAgentRequest struct {
	Description           string `json:"description"`
	GatherCycle           int32  `json:"gather_cycle"`
	ReportCycle           int32  `json:"report_cycle"`
	DataCompressAlgorithm string `json:"data_compress_algorithm"`
	EnableStreamAbility   bool   `json:"enable_stream_ability"`

	OpenApiDoc string `json:"open_api_doc"` // base64 encode
	EnableConf string `json:"enable_conf"`  // base64 encode
}
type PushTaskRequest struct {
	AgentID string `json:"agent_id"`
	TaskID  string `json:"task_id"`

	Task string `json:"task"` // base64 encode
}
