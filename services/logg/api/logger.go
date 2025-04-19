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
	"github.com/google/uuid"
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
	_root string
}

func (l *Logger) WithVersionJson(v []byte) Log {
	newLogger := *l
	newLogger.version = v
	return &newLogger
}
func (l *Logger) WithException(e *exception.Exception) Log {
	newLogger := *l
	newLogger.Exception = e
	return &newLogger
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
func (l *Logger) AsRoot() Log {
	newLogger := *l
	newLogger._root = uuid.NewString() // create a new root object
	fmt.Println(l.String())
	return &newLogger
}
func (l *Logger) NewLog() Log {
	return l.AsRoot()
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
	newLogger := *l
	newLogger.ctx = ctx
	return &newLogger
}
func (l *Logger) WithWaitOption(waitOption bool) Log {
	newLogger := *l
	newLogger.waitOption = waitOption
	return &newLogger
}
func (l *Logger) WithLoglevel(level svcpb.Level) Log {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *Logger) WithIP(ip string) Log {
	newLogger := *l
	newLogger.ip = ip
	return &newLogger
}

func (l *Logger) WithOperatorID(id string) Log {
	newLogger := *l
	newLogger.operatorID = id
	return &newLogger
}

func (l *Logger) WithProtocol(protocol string) Log {
	newLogger := *l
	newLogger.protocol = protocol
	return &newLogger
}

func (l *Logger) WithAction(action string) Log {
	newLogger := *l
	newLogger.action = action
	return &newLogger
}

func (l *Logger) WithOperationType(operationType svcpb.Operation) Log {
	newLogger := *l
	newLogger.operationType = operationType
	return &newLogger
}

func (l *Logger) WithCost(cost time.Duration) Log {
	newLogger := *l
	newLogger.cost = cost
	return &newLogger
}

func (l *Logger) WithMetaData(metadata any) Log {
	newLogger := *l
	newLogger.metadata = metadata
	return &newLogger
}

func (l *Logger) WithPrintToStdout() Log {
	newLogger := *l
	newLogger.printToStdout = true
	return &newLogger
}
func (l *Logger) WithNotPrintToStdout() Log {
	newLogger := *l
	newLogger.printToStdout = false
	return &newLogger
}

func (l *Logger) WithExceptionCode(code exceptionCode.ResCode) Log {
	newLogger := *l
	newLogger.exceptionCode = code
	return &newLogger
}

func (l *Logger) Debug(args ...interface{}) {
	l.message = fmt.Sprintf("%v", args)
	if l.printToStdout {
		logger.Debug(l.message)
	}
	//_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Info(args ...interface{}) {
	l.message = fmt.Sprintf("%v", args)
	if l.printToStdout {
		logger.Info(l.message)
	}
	if l.logLevel == svcpb.Level_LogUnknown {
		l.logLevel = svcpb.Level_INFO
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Warn(args ...interface{}) {
	l.message = fmt.Sprintf("%v", args)
	if l.printToStdout {
		logger.Warn(l.message)
	}
	if l.logLevel == svcpb.Level_LogUnknown {
		l.logLevel = svcpb.Level_WARN
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Error(args ...interface{}) {
	l.message = fmt.Sprintf("%v", args)
	if l.printToStdout {
		logger.Error(l.message)
	}
	if l.logLevel == svcpb.Level_LogUnknown {
		l.logLevel = svcpb.Level_ERROR
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.message = fmt.Sprintf(format, args...)
	if l.printToStdout {
		logger.Debug(l.message)
	}
	//_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.message = fmt.Sprintf(format, args...)
	if l.printToStdout {
		logger.Info(l.message)
	}
	if l.logLevel == svcpb.Level_LogUnknown {
		l.logLevel = svcpb.Level_INFO
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.message = fmt.Sprintf(format, args...)
	if l.printToStdout {
		logger.Warn(l.message)
	}
	if l.logLevel == svcpb.Level_LogUnknown {
		l.logLevel = svcpb.Level_WARN
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.message = fmt.Sprintf(format, args...)
	if l.printToStdout {
		logger.Error(l.message)
	}
	if l.logLevel == svcpb.Level_LogUnknown {
		l.logLevel = svcpb.Level_ERROR
	}
	_ = l.Transfer(l.ctx, l.generateProtobuf(), l.waitOption)
}
