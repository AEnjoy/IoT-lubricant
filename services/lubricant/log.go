package lubricant

import (
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
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
