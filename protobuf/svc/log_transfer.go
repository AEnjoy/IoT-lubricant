package svc

import "context"

type LogTransfer interface {
	// Transfer :transfer logs to core
	// if wait is true, this call will block until the request is completed and return the result(retval error)
	// if wait is false, this call will return immediately and the result(retval) will be nil
	Transfer(ctx context.Context, data *Logs, wait bool) (retval error)
}
