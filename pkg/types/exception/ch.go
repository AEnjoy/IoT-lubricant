package exception

import "github.com/AEnjoy/IoT-lubricant/pkg/logger"

var ErrCh = make(chan error)

func init() {
	go Handle()
}

func Handle() {
	for err := range ErrCh {
		logger.Errorln(err)
	}
}
