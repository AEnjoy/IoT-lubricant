package request

import "github.com/aenjoy/iot-lubricant/pkg/model"

type AgentsBindProjectRequest struct {
	ProjectID string   `json:"project_id"`
	AgentIDs  []string `json:"agent_ids"`
}
type AgentsBindWasherRequest struct {
	ProjectID string   `json:"project_id"`
	WasherID  int      `json:"washer_id"`
	AgentIDs  []string `json:"agent_ids"`
}
type AddWasherRequest struct {
	ProjectID   string   `json:"project_id"`
	Description string   `json:"description"`
	Table       *string  `json:"table,omitempty"`
	Interpreter string   `json:"interpreter"`
	Script      string   `json:"script"`
	Command     string   `json:"command"`
	Location    string   `json:"location,omitempty"` // core, gateway
	ToAgents    []string `json:"to_agents,omitempty"`
}
type AddProjectRequest struct {
	ProjectName string `json:"project_name"`
	Description string `json:"description"`

	DataBaseType  string            `json:"data_base_type,omitempty"` // mysql, TDEngine, mongodb
	StoreTable    string            `json:"store_table,omitempty"`
	DSNLinkerInfo *model.LinkerInfo `json:"dsn_linker_info,omitempty"`
	//Agents        []string          `json:"agents,omitempty"`
	Washer *AddWasherRequest `json:"washer,omitempty"`
}
