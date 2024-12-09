package exception

import (
	"errors"

	"github.com/AEnjoy/IoT-lubricant/pkg/types/code"
)

type Exception struct {
	Code         code.ResCode `json:"code"`
	Msg          string       `json:"msg"`
	Level        code.Level   `json:"level,omitempty"`
	Reason       interface{}  `json:"reason,omitempty"`
	DetailReason interface{}  `json:"detail_reason,omitempty"`
	Data         interface{}  `json:"data,omitempty"`
	Operation    Operation    `json:"-"`
	doOperation  bool
}
type Option func(*Exception)

func (e *Exception) Error() string {
	return e.Msg
}

// WithMsg 允许设置多条错误信息，但新旧错误信息间没有分隔符
func WithMsg(msg string) Option {
	return func(e *Exception) {
		e.Msg += msg
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
func New(code code.ResCode, opts ...Option) *Exception {
	exception := &Exception{
		Code: code,
	}

	for _, opt := range opts {
		opt(exception)
	}

	if exception.Operation != nil && exception.doOperation {
		exception.Operation.Do(exception)
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
