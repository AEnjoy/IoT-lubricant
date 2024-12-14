package core

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
)

var errCh = make(chan *exception.ErrLogInfo, 3)

func handleErrLog() {
	for e := range errCh {
		// todo: need to report error to user
		logger.Error(e.String())
	}
}
func init() {
	go handleErrLog()
}
