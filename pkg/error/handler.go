package error

import (
	"fmt"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
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
	return len(err.Exception.Msg) == 0
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
				Msg:    []string{message},
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
			for _, f := range globalErrorDo {
				f(&err)
			}
			return
		}
	default:
	}
	for _, f := range globalSuccessDo {
		f()
	}
	eh.successRun()
}

var globalSuccessDo []func()
var globalErrorDo []func(error)

func ErrorHandlerSetGlobalSuccessWillDo(callback func()) {
	globalSuccessDo = append(globalSuccessDo, callback)
}
func ErrorHandlerSetGlobalErrorWillDo(callback func(error)) {
	globalErrorDo = append(globalErrorDo, callback)
}
