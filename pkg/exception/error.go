package exception

import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
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
	Err error
	Msg string
}

func (err Error) IsEmpty() bool {
	return err.Err == nil
}

func NewErrorChan() *ErrorChan {
	return &ErrorChan{
		ErrCh: make(chan Error, 1),
	}
}

// Report sends error to the error channel
func (ec *ErrorChan) Report(err error, format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	if err != nil {
		logger.Debugf("Error reported: ", err, "Message: ", message)
		select {
		case ec.ErrCh <- Error{
			Err: err,
			Msg: message,
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
		logger.Errorf(errorMessageTemplate, err.Err.Error(), err.Msg)
		if !err.IsEmpty() {
			eh.errorRun(err.Err)
			return
		}
	default:
	}
	eh.successRun()
}
