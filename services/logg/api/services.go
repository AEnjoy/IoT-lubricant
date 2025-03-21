package api

import (
	"time"

	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
)

type Interface interface {
	WithLoglevel(level svcpb.Level) Interface
	WithIP(ip string) Interface
	// WithOperatorID : DeviceID / UserID
	WithOperatorID(id string) Interface
	WithProtocol(protocol string) Interface
	WithAction(action string) Interface
	WithOperationType(operationType svcpb.Operation) Interface
	WithCost(cost time.Duration) Interface
	// WithMetaData : metadata is data that can be serialized as JSON
	WithMetaData(metadata any) Interface
	WithPrintToStdout() Interface
	WithExceptionCode(code exceptionCode.ResCode) Interface

	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
