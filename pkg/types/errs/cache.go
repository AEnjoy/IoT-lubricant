package errs

import (
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
)

var (
	ErrNeedInit  error = exception.New(code.ErrorCacheNeedInit)
	ErrNullCache error = exception.New(code.ErrorCacheNullCache)
)
