package task

import "github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"

type ContainerDeployInfo struct {
	Name  string            `json:"name"`  // 容器名
	Image string            `json:"image"` // 容器镜像
	Port  int               `json:"port"`  // 暴露端口
	Env   map[string]string `json:"env"`   // 环境变量
}
type ContainerInfo struct {
	IP       string `json:"ip"`        // 容器内 Driver的IP
	HostName string `json:"host_name"` // optional
}
type ScheduleAddRequest struct {
	AgentID  string `json:"agent_id"`
	Interval int    `json:"interval"`
	Request  any    `json:"request"` // 请求体(POST)或请求参数(GET)
	Api      string `json:"api"`
}
type ScheduleAddResponse struct {
	ID string `json:"id"` // 计划任务ID
}
type ScheduleRemoveRequest struct {
	ID string `json:"id"` // 计划任务ID
}
type OpenAPIEnableRequest struct {
	AgentID       string          `json:"agent_id"`
	OpenAPIEnable openapi.ApiInfo `json:"openapi_enable"`
}
type OpenAPIDisableRequest struct {
	AgentID        string   `json:"agent_id"`
	OpenAPIDisable []string `json:"openapi_disable"` // api-paths
}
type SendRequest struct {
	AgentID string `json:"agent_id"`
	Method  string `json:"method"`
	Path    string `json:"path"`
	Request any    `json:"request"` // 请求体(POST)或请求参数(GET)
}
type OpenAPIGetRequest struct {
	AgentID string `json:"agent_id"`
}
type OpenAPIEnableGetRequest struct {
	AgentID string `json:"agent_id"`
}
