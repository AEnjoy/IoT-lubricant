package errs

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
)

var (
	ErrNeedTxn error = exception.New(code.ErrorDbNeedTxn)
)
