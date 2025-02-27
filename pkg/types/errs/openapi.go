package errs

import (
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
)

var (
	ErrNotFound      error = exception.New(exceptionCode.ErrorApiNotFound)
	ErrInvalidMethod error = exception.New(exceptionCode.ErrorApiInvalidMethod)

	ErrInvalidInput error = exception.New(exceptionCode.ErrorApiInvalidInput)
	ErrInvalidPath  error = exception.New(exceptionCode.ErrorApiInvalidPath)
	ErrInvalidSlot  error = exception.New(exceptionCode.ErrorApiInvalidSlot)
	ErrNotInit      error = exception.New(exceptionCode.ErrorApiNotInit)
)
