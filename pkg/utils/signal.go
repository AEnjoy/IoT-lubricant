package utils

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
)

func HandelExitSignal(ctxCancel, call, timeOutCall func(), timeout time.Duration) func() {
	return func() {
		// Handle SIGINT and SIGTERM
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		logger.Infoln("Received shutdown signal")
		if ctxCancel != nil {
			ctxCancel()
		}
		if call != nil {
			call()
		}
		if timeout != 0 && timeOutCall != nil {
			go func() {
				<-time.After(timeout)
				logger.Errorln("Shutdown timeout, force exit")
				timeOutCall()
				os.Exit(1)
			}()
		}
		os.Exit(0)
	}
}
