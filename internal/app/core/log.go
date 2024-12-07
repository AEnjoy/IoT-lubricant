package core

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/exception"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
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
