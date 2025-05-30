package gateway

import (
	level "github.com/aenjoy/iot-lubricant/pkg/types/code"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	object "github.com/aenjoy/iot-lubricant/pkg/types/task"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *app) _handelAgentControlError(ch <-chan *exception.Exception) {
	makeMessage := func(e *exception.Exception) *corepb.ReportRequest {
		var agentid string
		if s := e.Get(string(object.TargetAgent)); s != nil {
			agentid = s.(string)
		}
		var retVal = corepb.ReportRequest{
			GatewayId: gatewayId,
			AgentId:   agentid,
			Req: &corepb.ReportRequest_Error{
				Error: &corepb.ReportErrorRequest{
					ErrorMessage: &metapb.ErrorMessage{
						ErrorType: func() *int32 {
							if e.Level != level.Unknown {
								var retVal = int32(e.Level)
								return &retVal
							}
							return nil
						}(),
						Code: &status.Status{
							Message: e.Error(),
							Code:    int32(e.Code),
						},
						Time: timestamppb.Now(),
						Module: func() *string {
							var retVal = "agent"
							return &retVal
						}(),
					},
				},
			},
		}
		return &retVal
	}
	for e := range ch {
		// todo: report error to gateway when level >= error
		if e != nil {
			switch e.Level {
			case level.Debug:
				logg.L.Debugf("%v", e)
			case level.Info:
				logg.L.Infof("%v", e)
			case level.Warn:
				logg.L.Warnf("%v", e)
				_reportMessage <- makeMessage(e)
			default:
				logg.L.Errorf("%v", e)
				_reportMessage <- makeMessage(e)
			}
		}
	}
}
