package errs

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidMethod = errors.New("invalid method")

	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidPath  = errors.New("invalid path")
	ErrInvalidSlot  = errors.New("invalid slot")
	ErrNotInit      = errors.New("not initialized")
)
