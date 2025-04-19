package exception

import "github.com/aenjoy/iot-lubricant/pkg/logger"

var ErrCh = make(chan error)

func init() {
	go Handle()
}

func Handle() {
	for err := range ErrCh {
		logger.Errorln(err)
	}
}
