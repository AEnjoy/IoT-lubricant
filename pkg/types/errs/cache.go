package errs

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
)

var (
	ErrNeedInit  error = exception.New(code.ErrorCacheNeedInit)
	ErrNullCache error = exception.New(code.ErrorCacheNullCache)
)
