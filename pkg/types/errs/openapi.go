package errs

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
)

var (
	ErrNotFound      error = exception.New(code.ErrorApiNotFound)
	ErrInvalidMethod error = exception.New(code.ErrorApiInvalidMethod)

	ErrInvalidInput error = exception.New(code.ErrorApiInvalidInput)
	ErrInvalidPath  error = exception.New(code.ErrorApiInvalidPath)
	ErrInvalidSlot  error = exception.New(code.ErrorApiInvalidSlot)
	ErrNotInit      error = exception.New(code.ErrorApiNotInit)
)
