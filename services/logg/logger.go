package logg

import (
	"fmt"
	"time"

	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
	"github.com/aenjoy/iot-lubricant/services/logg/api"
)

type Logger struct {
	logLevel      svcpb.Level
	ip            string
	operatorID    string
	protocol      string
	action        string
	operationType svcpb.Operation
	cost          time.Duration
	metadata      any
	printToStdout bool
	exceptionCode exceptionCode.ResCode
}

func (l *Logger) WithLoglevel(level svcpb.Level) api.Interface {
	l.logLevel = level
	return &*l
}

func (l *Logger) WithIP(ip string) api.Interface {
	l.ip = ip
	return &*l
}

func (l *Logger) WithOperatorID(id string) api.Interface {
	l.operatorID = id
	return &*l
}

func (l *Logger) WithProtocol(protocol string) api.Interface {
	l.protocol = protocol
	return &*l
}

func (l *Logger) WithAction(action string) api.Interface {
	l.action = action
	return &*l
}

func (l *Logger) WithOperationType(operationType svcpb.Operation) api.Interface {
	l.operationType = operationType
	return &*l
}

func (l *Logger) WithCost(cost time.Duration) api.Interface {
	l.cost = cost
	return &*l
}

func (l *Logger) WithMetaData(metadata any) api.Interface {
	l.metadata = metadata
	return &*l
}

func (l *Logger) WithPrintToStdout() api.Interface {
	l.printToStdout = true
	return &*l
}

func (l *Logger) WithExceptionCode(code exceptionCode.ResCode) api.Interface {
	l.exceptionCode = code
	return &*l
}

func (l *Logger) Debug(args ...interface{}) {
	if l.printToStdout {
		fmt.Println("[DEBUG]", args)
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.printToStdout {
		fmt.Println("[INFO]", args)
	}
}

func (l *Logger) Warn(args ...interface{}) {
	if l.printToStdout {
		fmt.Println("[WARN]", args)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.printToStdout {
		fmt.Println("[ERROR]", args)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.printToStdout {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.printToStdout {
		fmt.Printf("[INFO] "+format+"\n", args...)
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.printToStdout {
		fmt.Printf("[WARN] "+format+"\n", args...)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.printToStdout {
		fmt.Printf("[ERROR] "+format+"\n", args...)
	}
}
