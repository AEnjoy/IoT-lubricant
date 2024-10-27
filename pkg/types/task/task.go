package task

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/types/user"
)

type Task struct {
	TaskID           string    `json:"id" gorm:"column:id;primary_key"`
	Operator         user.Role `json:"operator" gorm:"column:operator;type:tiny"`
	Executor         user.Role `json:"executor" gorm:"column:executor;type:tiny"`
	ExecutorID       string    `json:"executor_id" gorm:"column:executor_id"` // user/core/gateway/agent ID UUID
	OperationType    Operation `json:"operation" gorm:"column:operation;type:tiny"`
	OperationCommend string    `json:"operation_commend" gorm:"column:operation_commend;serializer:json"` //json
	SupportRollback  bool      `json:"support_rollback" gorm:"column:support_rollback"`

	OperationTime int64 `json:"operation_time" gorm:"column:created_at"`
	UpdatedAt     int64 `json:"-" gorm:"column:updated_at"`
}

func (Task) TableName() string {
	return "task_log"
}
