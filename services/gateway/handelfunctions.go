package gateway

import (
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	level "github.com/aenjoy/iot-lubricant/pkg/types/code"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
)

func HandelAgentControlError(ch <-chan *exception.Exception) {
	for e := range ch {
		// todo: report error to gateway when level >= error
		if e != nil {
			switch e.Level {
			case level.Debug:
				logger.Debugf("%v", e)
			case level.Info:
				logger.Infof("%v", e)
			case level.Warn:
				logger.Warnf("%v", e)
			default:
				logger.Errorf("%v", e)
			}
		}
	}
}
