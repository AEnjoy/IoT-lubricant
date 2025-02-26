package errs

import (
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
)

var (
	ErrNeedTxn error = exception.New(code.ErrorDbNeedTxn)
)
