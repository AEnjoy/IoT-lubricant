package exception

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/code"
)

var ErrCh = make(chan error)

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
func init() {
	go Handle()
}

func Handle() {
	for err := range ErrCh {
		logger.Errorln(err)
	}
}
