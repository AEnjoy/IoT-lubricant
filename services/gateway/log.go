package gateway

import (
	"context"

	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
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
	return nil
}
func (a *app) grpcReportLogApp() {
	for logs := range logCollect {
		_reportMessage <- &corepb.ReportRequest{
			GatewayId: gatewayId,
			Req: &corepb.ReportRequest_ReportLog{
				ReportLog: &corepb.ReportLogRequest{
					GatewayId: gatewayId,
					Logs:      []*svcpb.Logs{logs},
				},
			},
		}
	}
}
