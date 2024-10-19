package app

import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
)

var errCh = make(chan *ErrLogInfo, 3)

type ErrLogInfo struct {
	User    string
	Agent   string
	Message error
}

func (r *ErrLogInfo) String() string {
	return fmt.Sprintf("User:%s 's Agent:%s report an error:%s ", r.User, r.Agent, r.Message.Error())
}

func handleErrLog() {
	for e := range errCh {
		// todo: need to report error to user
		logger.Error(e.String())
	}
}
func init() {
	go handleErrLog()
}
