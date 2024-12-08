package exception

import "github.com/AEnjoy/IoT-lubricant/pkg/types/code"

type Exception struct {
	Code         code.ResCode `json:"code"`
	Msg          string       `json:"msg"`
	Reason       interface{}  `json:"reason,omitempty"`
	DetailReason interface{}  `json:"detail_reason,omitempty"`
	Data         interface{}  `json:"data,omitempty"`
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

func New(code code.ResCode, opts ...Option) *Exception {
	exception := &Exception{
		Code: code,
	}

	for _, opt := range opts {
		opt(exception)
	}

	return exception
}
