package request

import "github.com/aenjoy/iot-lubricant/pkg/types/container"

type AddAgentRequest struct {
	Description           string `json:"description"`
	GatherCycle           int32  `json:"gather_cycle"`
	ReportCycle           int32  `json:"report_cycle"`
	Address               string `json:"address"` // ip:port optional
	DataCompressAlgorithm string `json:"data_compress_algorithm"`
	EnableStreamAbility   bool   `json:"enable_stream_ability"`

	AgentContainerInfo  *container.Container `json:"agent_container_info,omitempty"`
	DriverContainerInfo *container.Container `json:"driver_container_info,omitempty"`

	OpenApiDoc string `json:"open_api_doc"` // base64 encode
	EnableConf string `json:"enable_conf"`  // base64 encode
}
type PushTaskRequest struct {
	AgentID string `json:"agent_id"`
	TaskID  string `json:"task_id"`

	Task string `json:"task"` // base64 encode
}
