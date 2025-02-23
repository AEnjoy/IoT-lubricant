package task

type Target string

const (
	TargetCore    Target = "lubricant"
	TargetGateway Target = "gateway"
	TargetAgent   Target = "agent"
)

type TaskType int

const (
	TaskType_Unknow TaskType = iota
	TaskType_StartAgentRequest
	TaskType_CreateAgentRequest
	TaskType_EditAgentRequest
	TaskType_RemoveAgentRequest
	TaskType_StopAgentRequest
	TaskType_UpdateAgentRequest
)
