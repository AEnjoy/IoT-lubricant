package task

import "github.com/AEnjoy/IoT-lubricant/pkg/types"

type Operation uint8

const (
	_                              Operation = 0
	OperationUserLogin             Operation = 10
	OperationUserLogout            Operation = 11
	OperationUserChangePassword    Operation = 12
	OperationAddTask               Operation = 13
	OperationQueryTask             Operation = 14
	OperationViewTaskResult        Operation = 15
	OperationAddGateway            Operation = 20
	OperationRemoveGateway         Operation = 21
	OperationAddAgent              Operation = 22
	OperationRemoveAgent           Operation = 23
	OperationAddSchedule           Operation = 24
	OperationRemoveSchedule        Operation = 25
	OperationAddDriverContainer    Operation = 30
	OperationRemoveDriverContainer Operation = 31
	OperationAddAgentContainer     Operation = 32
	OperationRemoveAgentContainer  Operation = 33
	OperationEnableOpenAPI         Operation = 40
	OperationDisableOpenAPI        Operation = 41
	OperationSendRequest           Operation = 42
	OperationGetOpenAPIDoc         Operation = 43
	OperationGetEnableOpenAPI      Operation = 44
)

type Task struct {
	TaskID           string     `json:"id" gorm:"column:id;primary_key"`
	Operator         types.Role `json:"operator" gorm:"column:operator;type:tiny"`
	Executor         types.Role `json:"executor" gorm:"column:executor;type:tiny"`
	ExecutorID       string     `json:"executor_id" gorm:"column:executor_id"` // user/core/gateway/agent ID UUID
	OperationType    Operation  `json:"operation" gorm:"column:operation;type:tiny"`
	OperationCommend string     `json:"operation_commend" gorm:"column:operation_commend;serializer:json"` //json
	SupportRollback  bool       `json:"support_rollback" gorm:"column:support_rollback"`

	OperationTime int64 `json:"operation_time" gorm:"column:created_at"`
	UpdatedAt     int64 `json:"-" gorm:"column:updated_at"`
}
