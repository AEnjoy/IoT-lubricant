package api

import (
	"context"
	"errors"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
)

type Log interface {
	// AsRoot : Use the currently set value as the root log branch
	AsRoot() Log
	// NewLog : It is an alias of AsRoot() Log
	NewLog() Log

	Reset()

	WithLoglevel(level svcpb.Level) Log
	WithVersionJson(v []byte) Log
	WithIP(ip string) Log
	// WithOperatorID : DeviceID / UserID
	WithOperatorID(id string) Log
	WithProtocol(protocol string) Log
	WithAction(action string) Log
	WithOperationType(operationType svcpb.Operation) Log
	WithCost(cost time.Duration) Log
	// WithMetaData : metadata is data that can be serialized as JSON
	WithMetaData(metadata any) Log
	WithPrintToStdout() Log
	WithException(e *exception.Exception) Log
	// WithExceptionCode : If WithException has already specified an exceptionCode, it is not necessary to call WithExceptionCode
	WithExceptionCode(code exceptionCode.ResCode) Log

	WithContext(ctx context.Context) Log
	WithWaitOption(waitOption bool) Log

	String() string

	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func NewLogger(transfer svcpb.LogTransfer, wait bool) (Log, error) {
	if transfer == nil {
		return nil, errors.New("transfer is nil")
	}
	return &Logger{LogTransfer: transfer, ctx: context.TODO(), waitOption: wait}, nil
}
