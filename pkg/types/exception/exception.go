package exception

import (
	"errors"
	"strings"

	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/code"
	except "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
)

type Exception struct {
	Code         except.ResCode `json:"code"`
	Msg          []string       `json:"msg"`
	Level        code.Level     `json:"level,omitempty"`
	Reason       interface{}    `json:"reason,omitempty"`
	DetailReason interface{}    `json:"detail_reason,omitempty"`
	Data         interface{}    `json:"data,omitempty"`
	Operation    Operation      `json:"-"`
	ShowLog      bool           `json:"-"`
	doOperation  bool
}
type Option func(*Exception)

func (e *Exception) Error() string {
	var str strings.Builder
	for _, msg := range e.Msg {
		str.WriteString(msg)
		str.WriteString("; ")
	}
	return str.String()
}

// WithMsg 允许设置多条错误信息，但新旧错误信息间没有分隔符
func WithMsg(msg string) Option {
	return func(e *Exception) {
		e.Msg = append(e.Msg, msg)
	}
}

func WithReason(reason interface{}) Option {
	return func(e *Exception) {
		e.Reason = reason
	}
}

func WithDetailReason(detailReason interface{}) Option {
	return func(e *Exception) {
		e.DetailReason = detailReason
	}
}

func WithData(data interface{}) Option {
	return func(e *Exception) {
		e.Data = data
	}
}

func WithLevel(l code.Level) Option {
	return func(e *Exception) {
		e.Level = l
	}
}
func WithWillDo(callback func()) Option {
	return func(e *Exception) {
		callback()
	}
}
func WithErrWillDo(callback func(error)) Option {
	return func(e *Exception) {
		callback(e)
	}
}
func WithOperation(operation Operation, do bool) Option {
	return func(e *Exception) {
		e.Operation = operation
		e.doOperation = do
	}
}
func New(c except.ResCode, opts ...Option) *Exception {
	exception := &Exception{
		Code: c,
	}

	m := c.GetMsg()
	if m != except.StatusMsgMap[except.ErrorUnknown] {
		exception.Msg = []string{m}
	}

	for _, opt := range opts {
		opt(exception)
	}

	if exception.Operation != nil && exception.doOperation {
		_ = exception.Operation.Do(exception)
	}
	if exception.ShowLog {
		logger.Errorf("code: %d, msg: %s, reason: %v, detail_reason: %v, data: %v",
			exception.Code, exception.Msg, exception.Reason, exception.DetailReason, exception.Data)
	}
	return exception
}
func ErrNewException(err error, code except.ResCode, opts ...Option) *Exception {
	if err == nil {
		return nil
	}
	var exception *Exception
	if errors.As(err, &exception) {
		exception.Code = code
		for _, opt := range opts {
			opt(exception)
		}
	} else {
		opts = append(opts, WithMsg(err.Error()))
		exception = New(code, opts...)
	}
	return exception
}
func CheckException(err error) (*Exception, error) {
	var e *Exception
	if ok := errors.As(err, &e); ok {
		return e, nil
	}
	return nil, errors.New("not an internal exception")
}
func WithLogShow() Option {
	return func(e *Exception) {
		e.ShowLog = true
	}
}
