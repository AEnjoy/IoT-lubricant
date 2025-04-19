package logg

import (
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func (a *app) handel(in any) {
	data := in.([]byte)
	pb := &svcpb.Logs{}
	err := proto.Unmarshal(data, pb)
	if err != nil {
		logger.Errorf("failed to unmarshal data: %v", err)
		return
	}
	err = a.db.Write(a.ctx, &model.Log{
		LogID:      uuid.NewString(),
		OperatorID: pb.OperatorID,

		ServiceName:   pb.ServiceName,
		Level:         pb.Level,
		IPAddress:     pb.IPAddress,
		Protocol:      pb.Protocol,
		Action:        pb.Action,
		OperationType: pb.OperationType,
		Cost:          pb.Cost.AsTime().Unix(),
		Message:       pb.Message,
		ServiceErrorCode: func() exceptionCode.ResCode {
			if pb.ServiceErrorCode == nil {
				return exceptionCode.EmptyValue
			}
			return exceptionCode.ResCode(*pb.ServiceErrorCode)
		}(),
		Metadata:      string(pb.Metadata),
		ExceptionInfo: string(pb.ExceptionInfo),
		Time:          pb.Time.AsTime(),
	})
	if err != nil {
		logger.Errorf("failed to write log: %v", err)
	}
}
