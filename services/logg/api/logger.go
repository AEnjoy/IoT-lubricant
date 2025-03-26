package api

import (
	"context"
	"fmt"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
	"github.com/bytedance/sonic"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Logger struct {
	logLevel   svcpb.Level
	ip         string
	operatorID string
	protocol   string
	action     string
	message    string
	version    []byte // json 格式
	*exception.Exception
	operationType svcpb.Operation
	cost          time.Duration
	metadata      any
	printToStdout bool
	exceptionCode exceptionCode.ResCode

	ctx        context.Context
	waitOption bool
	svcpb.LogTransfer
}

func (l *Logger) WithVersionJson(v []byte) Log {
	l.version = v
	return l
}
func (l *Logger) WithException(e *exception.Exception) Log {
	l.Exception = e
	return l
}
func (l *Logger) Reset() {
	l.logLevel = svcpb.Level_DEBUG
	l.ip = ""
	l.operatorID = ""
	l.protocol = ""
	l.action = ""
	l.message = ""
	l.version = nil
	l.Exception = nil
	l.operationType = svcpb.Operation_Unknown
	l.cost = 0
	l.metadata = nil
	l.printToStdout = false
	l.exceptionCode = exceptionCode.EmptyValue
}
func (l *Logger) Root() Log {
	return &*l
}
func (l *Logger) String() string {
	return l.generateProtobuf().String()
}
func (l *Logger) generateProtobuf() *svcpb.Logs {
	return &svcpb.Logs{
		ID:            "",
		Time:          timestamppb.New(time.Now()),
		ServiceName:   ServiceName,
		Level:         l.logLevel,
		IPAddress:     l.ip,
		Action:        l.action,
		Protocol:      l.protocol,
		OperationType: l.operationType,
		OperatorID:    l.operatorID,
		Cost: func() *timestamppb.Timestamp {
			if l.cost != 0 {
				return timestamppb.New(time.Now().Add(-l.cost))
			}
			return nil
		}(),
		Message: l.message,
		Version: l.version,
		ServiceErrorCode: func() *int32 {
			if l.exceptionCode != exceptionCode.ErrorUnknown && l.exceptionCode != exceptionCode.EmptyValue {
				var retVal = int32(l.exceptionCode)
				return &retVal
			}
			return nil
		}(),
		ExceptionInfo: func() (retval []byte) {
			if l.Exception != nil {
				retval, _ = sonic.Marshal(l.Exception)
			}
			return
		}(),
		Metadata: func() (retval []byte) {
			if l.metadata != nil {
				retval, _ = sonic.Marshal(l.metadata)
			}
			return
		}(),
	}
}
func (l *Logger) WithContext(ctx context.Context) Log {
	l.ctx = ctx
	return l
}
func (l *Logger) WithWaitOption(waitOption bool) Log {
	l.waitOption = waitOption
	return l
}
func (l *Logger) WithLoglevel(level svcpb.Level) Log {
	l.logLevel = level
	return l
}

func (l *Logger) WithIP(ip string) Log {
	l.ip = ip
	return &*l
}

func (l *Logger) WithOperatorID(id string) Log {
	l.operatorID = id
	return &*l
}

func (l *Logger) WithProtocol(protocol string) Log {
	l.protocol = protocol
	return &*l
}

func (l *Logger) WithAction(action string) Log {
	l.action = action
	return &*l
}

func (l *Logger) WithOperationType(operationType svcpb.Operation) Log {
	l.operationType = operationType
	return &*l
}

func (l *Logger) WithCost(cost time.Duration) Log {
	l.cost = cost
	return &*l
}

func (l *Logger) WithMetaData(metadata any) Log {
	l.metadata = metadata
	return &*l
}

func (l *Logger) WithPrintToStdout() Log {
	l.printToStdout = true
	return &*l
}

func (l *Logger) WithExceptionCode(code exceptionCode.ResCode) Log {
	l.exceptionCode = code
	return &*l
}

func (l *Logger) Debug(args ...interface{}) {
	l.message = fmt.Sprintf("%v", args)
	if l.printToStdout {
		logger.Debug(l.message)
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Info(args ...interface{}) {
	l.message = fmt.Sprintf("%v", args)
	if l.printToStdout {
		logger.Info(l.message)
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Warn(args ...interface{}) {
	l.message = fmt.Sprintf("%v", args)
	if l.printToStdout {
		logger.Warn(l.message)
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Error(args ...interface{}) {
	l.message = fmt.Sprintf("%v", args)
	if l.printToStdout {
		logger.Error(l.message)
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.message = fmt.Sprintf(format, args)
	if l.printToStdout {
		logger.Debug(l.message)
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.message = fmt.Sprintf(format, args)
	if l.printToStdout {
		logger.Info(l.message)
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.message = fmt.Sprintf(format, args)
	if l.printToStdout {
		logger.Warn(l.message)
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.message = fmt.Sprintf(format, args)
	if l.printToStdout {
		logger.Error(l.message)
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}
