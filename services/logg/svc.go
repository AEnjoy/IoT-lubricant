package logg

import (
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
	"google.golang.org/protobuf/proto"
)

func (a *app) handel(data []byte) {
	pb := &svcpb.Logs{}
	err := proto.Unmarshal(data, pb)
	if err != nil {
		logger.Errorf("failed to unmarshal data: %v", err)
		return
	}
	_ = a.db.Write(a.ctx, &model.Log{})
}
