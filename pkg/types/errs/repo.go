package errs

import (
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
)

var (
	ErrNeedTxn error = exception.New(exceptionCode.ErrorDbNeedTxn)
)
