package error

import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/code"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
)

const errorMessageTemplate = "error message: %s, Message: %s"

type ErrorHandler struct {
	ErrorChan  ErrorChan
	successRun func()
	errorRun   func(error)
}

type ErrorChan struct {
	ErrCh chan Error
}

type Error struct {
	exception.Exception
}

func (err Error) IsEmpty() bool {
	return err.Exception.Msg == ""
}

func NewErrorChan() *ErrorChan {
	return &ErrorChan{
		ErrCh: make(chan Error, 1),
	}
}

// Report sends error to the error channel
func (ec *ErrorChan) Report(err error, code code.ResCode, format string, useLogger bool, a ...any) {
	message := fmt.Sprintf(format, a...)
	if err != nil {
		if useLogger {
			logger.Errorf(errorMessageTemplate, message, err)
		}
		select {
		case ec.ErrCh <- Error{
			Exception: exception.Exception{
				Code:   code,
				Msg:    message,
				Reason: err,
			},
		}:
		default:
		}
	}
}

func HandleErrorCh(ec *ErrorChan) *ErrorHandler {
	var errorHandler = &ErrorHandler{
		ErrorChan:  *ec,
		successRun: func() {},
		errorRun:   func(err error) {},
	}

	return errorHandler
}

func (eh *ErrorHandler) ErrorWillDo(callback func(error)) *ErrorHandler {
	eh.errorRun = callback
	return eh
}

func (eh *ErrorHandler) SuccessWillDo(callback func()) *ErrorHandler {
	eh.successRun = callback
	return eh
}

func (eh *ErrorHandler) Do() {
	select {
	case err := <-eh.ErrorChan.ErrCh:
		if !err.IsEmpty() {
			eh.errorRun(&err.Exception)
			return
		}
	default:
	}
	eh.successRun()
}
