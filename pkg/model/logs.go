package model

import (
	"time"

	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
)

type Log struct {
	ID         int    `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	LogID      string `json:"logId" gorm:"column:log_id;index"`
	OperatorID string `json:"operatorId" gorm:"column:operator_id;index"` // UserID / DevicesID

	ServiceName      string                `json:"serviceName" gorm:"column:service_name;index"`
	Version          string                `json:"version" gorm:"column:version;type:json"`
	Level            svcpb.Level           `json:"level" gorm:"column:level;type:tinyint;index"`
	IPAddress        string                `json:"ipAddress" gorm:"column:ip_address"`
	Protocol         string                `json:"protocol" gorm:"column:protocol"`
	Action           string                `json:"action" gorm:"column:action;index"`
	OperationType    svcpb.Operation       `json:"operationType" gorm:"column:operation_type;type:tinyint;index"`
	Cost             int64                 `json:"cost" gorm:"column:cost;type:timestamp"`
	Message          string                `json:"message" gorm:"column:message"`
	ServiceErrorCode exceptionCode.ResCode `json:"serviceErrorCode" gorm:"column:service_error_code;type:int;index"`
	Metadata         string                `json:"metadata" gorm:"column:metadata;type:json"`
	ExceptionInfo    string                `json:"exceptionInfo" gorm:"column:exception_info;type:json"`
	Time             time.Time             `json:"time" gorm:"column:time;type:timestamp"`
}

func (Log) TableName() string {
	return "logs"
}
