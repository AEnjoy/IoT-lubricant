package logg

import (
	"database/sql"

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
		ServiceName: func() sql.NullString {
			if len(pb.ServiceName) == 0 {
				return sql.NullString{}
			}
			return sql.NullString{String: pb.ServiceName, Valid: true}
		}(),
		Version: func() sql.NullString {
			if len(pb.Version) == 0 {
				return sql.NullString{}
			}
			return sql.NullString{String: string(pb.Version), Valid: true}
		}(),
		Level: pb.Level,
		IPAddress: func() sql.NullString {
			if len(pb.IPAddress) == 0 {
				return sql.NullString{}
			}
			return sql.NullString{String: pb.IPAddress, Valid: true}
		}(),
		Protocol: func() sql.NullString {
			if len(pb.Protocol) == 0 {
				return sql.NullString{}
			}
			return sql.NullString{String: pb.Protocol, Valid: true}
		}(),
		Action: func() sql.NullString {
			if len(pb.Action) == 0 {
				return sql.NullString{}
			}
			return sql.NullString{String: pb.Action, Valid: true}
		}(),
		OperationType: pb.OperationType,
		Cost: func() sql.NullInt64 {
			if pb.Time == nil {
				return sql.NullInt64{}
			}
			if !pb.Time.IsValid() {
				return sql.NullInt64{}
			}
			return sql.NullInt64{Int64: pb.Cost.AsTime().Unix(), Valid: true}
		}(),
		Message: func() sql.NullString {
			if len(pb.Message) == 0 {
				return sql.NullString{}
			}
			return sql.NullString{String: pb.Message, Valid: true}
		}(),
		ServiceErrorCode: func() exceptionCode.ResCode {
			if pb.ServiceErrorCode == nil {
				return exceptionCode.EmptyValue
			}
			return exceptionCode.ResCode(*pb.ServiceErrorCode)
		}(),
		Metadata: func() sql.NullString {
			if len(pb.Metadata) == 0 {
				return sql.NullString{}
			}
			return sql.NullString{String: string(pb.Metadata), Valid: true}
		}(),
		ExceptionInfo: func() sql.NullString {
			if len(pb.ExceptionInfo) == 0 {
				return sql.NullString{}
			}
			return sql.NullString{String: string(pb.ExceptionInfo), Valid: true}
		}(),
		Time: func() sql.NullTime {
			if pb.Time == nil {
				return sql.NullTime{}
			}
			if !pb.Time.IsValid() {
				return sql.NullTime{}
			}
			return sql.NullTime{Time: pb.Time.AsTime(), Valid: true}
		}(),
	})
	if err != nil {
		logger.Errorf("failed to write log: %v", err)
	}
}
