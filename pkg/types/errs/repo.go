package errs

import "errors"

var (
	ErrNeedTxn = errors.New("this operation need start with txn support")
)
