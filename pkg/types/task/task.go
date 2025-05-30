package task

import (
	"database/sql"

	"github.com/aenjoy/iot-lubricant/pkg/types/operation"
	"github.com/aenjoy/iot-lubricant/pkg/types/user"
)

type Task struct {
	TaskID           string              `json:"id" gorm:"column:id;primary_key"`
	Operator         user.Role           `json:"operator" gorm:"column:operator;type:TINYINT"`
	Executor         user.Role           `json:"executor" gorm:"column:executor;type:TINYINT"`
	ExecutorID       string              `json:"executor_id" gorm:"column:executor_id"` // user/core/gateway/agent ID UUID
	OperationType    operation.Operation `json:"operation" gorm:"column:operation;type:TINYINT"`
	OperationCommend string              `json:"operation_commend" gorm:"column:operation_commend;serializer:json"` //json
	SupportRollback  bool                `json:"support_rollback" gorm:"column:support_rollback"`

	OperationTime int64        `json:"operation_time" gorm:"column:created_at"`
	UpdatedAt     sql.NullTime `json:"-" gorm:"column:updated_at"`
}

func (Task) TableName() string {
	return "task_log"
}

const (
	TaskStatusUnknow int32 = iota
	TaskStatusPending
)
