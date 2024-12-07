package exception

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
)

var ErrCh = make(chan error)

type Exception struct {
	Code int
	Msg  string
}

func (e *Exception) Error() string {
	return e.Msg
}
func New(code int, msg ...string) *Exception {
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
