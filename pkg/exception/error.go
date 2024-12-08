package exception

import "github.com/AEnjoy/IoT-lubricant/pkg/types/code"

type Exception struct {
	Code         code.ResCode `json:"code"`
	Msg          string       `json:"msg"`
	Reason       interface{}  `json:"reason,omitempty"`
	DetailReason interface{}  `json:"detail_reason,omitempty"`
	Data         interface{}  `json:"data,omitempty"`
}

func (e *Exception) Error() string {
	return e.Msg
}
func New(code code.ResCode, msg ...string) *Exception {
	var msgs string
	for _, t := range msg {
		msgs += t
	}
	return &Exception{
		Code: code,
		Msg:  msgs,
	}
}
