package request

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
	ToAgents    []string `json:"to_agents,omitempty"`
}
