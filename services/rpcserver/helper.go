package rpcserver

import (
	"context"
	"fmt"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/types"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func (*PbCoreServiceImpl) taskSendErrorMessage(s grpc.BidiStreamingServer[corepb.Task, corepb.Task], code int, message string) {
	var out corepb.Task
	var errorResp corepb.Task_ErrorMessage
	errorResp.ErrorMessage = &metapb.ErrorMessage{Code: &status.Status{Code: int32(code), Message: message}}
	out.Task = &errorResp
	_ = s.Send(&out)
}
func (*PbCoreServiceImpl) gatewayOffline(mq mq.Mq, userid, gatewayid string) error {
	logg.L.Debugf("gateway offline: %s", gatewayid)
	return mq.PublishBytes(fmt.Sprintf("/monitor/%s/offline", taskTypes.TargetGateway),
		[]byte(fmt.Sprintf("%s<!SPLIT!>%s", userid, gatewayid)))
}

// getGatewayID 获取调用接口的网关ID
//
//	todo:(现阶段其实应该叫做网关名,后续需要使用真正的id)
func (*PbCoreServiceImpl) getGatewayID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata not found")
	}

	gatewayIDs := md.Get(string(types.NameGatewayID))
	if len(gatewayIDs) == 0 {
		return "", fmt.Errorf("gatewayid not found in metadata")
	}

	return gatewayIDs[0], nil
}
func (*PbCoreServiceImpl) getUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata not found")
	}

	userIDs := md.Get(constant.USER_ID)
	if len(userIDs) == 0 {
		return "", fmt.Errorf("userid not found in metadata")
	}

	return userIDs[0], nil
}

func (i *PbCoreServiceImpl) handelRecvData(data *corepb.Data, userId string) {
	err := i.pool.Submit(func() {
		dataCli := i.DataStore
		marshal, err := proto.Marshal(data)
		if err != nil {
			logg.L.Errorf("failed to marshal data: %v", err)
		}
		err = dataCli.V2mq.Publish(constant.DATASTORE_USER, []byte(userId))
		if err != nil {
			logg.L.Errorf("failed to publish userid[%s] to datastore topic: %v", userId, err)
			return
		}

		err = dataCli.V2mq.Publish(fmt.Sprintf(constant.DATASTORE_USER_DATA, userId), marshal)
		if err != nil {
			logg.L.Errorf("failed to publish data: %v", err)
		}
	})
	if err != nil {
		logg.L.Errorf("failed to create save data to store task thread")
	}
}
