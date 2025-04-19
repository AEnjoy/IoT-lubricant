package model

import (
	"database/sql"

	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
)

type Log struct {
	ID         int    `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	LogID      string `json:"logId" gorm:"column:log_id;index"`
	OperatorID string `json:"operatorId" gorm:"column:operator_id;index"` // UserID / DevicesID

	ServiceName      sql.NullString        `json:"serviceName" gorm:"column:service_name;index"`
	Version          sql.NullString        `json:"version" gorm:"column:version;type:json"`
	Level            svcpb.Level           `json:"level" gorm:"column:level;type:tinyint;index"`
	IPAddress        sql.NullString        `json:"ipAddress" gorm:"column:ip_address"`
	Protocol         sql.NullString        `json:"protocol" gorm:"column:protocol"`
	Action           sql.NullString        `json:"action" gorm:"column:action;index"`
	OperationType    svcpb.Operation       `json:"operationType" gorm:"column:operation_type;type:tinyint;index"`
	Cost             sql.NullInt64         `json:"cost" gorm:"column:cost"`
	Message          sql.NullString        `json:"message" gorm:"column:message"`
	ServiceErrorCode exceptionCode.ResCode `json:"serviceErrorCode" gorm:"column:service_error_code;type:int;index"`
	Metadata         sql.NullString        `json:"metadata" gorm:"column:metadata;type:json"`
	ExceptionInfo    sql.NullString        `json:"exceptionInfo" gorm:"column:exception_info;type:json"`
	Time             sql.NullTime          `json:"time" gorm:"column:time;type:timestamp"`
}

func (Log) TableName() string {
	return "logs"
}
