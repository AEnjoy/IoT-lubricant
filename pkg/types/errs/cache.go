package errs

import "errors"

var (
	ErrNeedInit  = errors.New("cache client need init")
	ErrNullCache = errors.New("cache client is nil")
)
