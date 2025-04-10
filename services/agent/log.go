package agent

import (
	"context"

	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
)

var logCollect = make(chan *svcpb.Logs, 200)

func (a *app) Transfer(ctx context.Context, data *svcpb.Logs, wait bool) (retval error) {
	//if wait{
	//	ctx,cancel:=context.WithTimeout(ctx,3 * time.Second)
	//	defer cancel()
	//	select {
	//	case logCollect <- data:
	//		return nil
	//	case <-ctx.Done():
	//		return ctx.Err()
	//	}
	//}
	go func() {
		logCollect <- data
	}()
	return
}
