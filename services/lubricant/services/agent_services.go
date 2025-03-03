package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/pkg/types/user"
	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/aenjoy/iot-lubricant/protobuf/gateway"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (a *AgentService) GetAgentStatus(ctx context.Context, taskid *string, gatewayid string, ids []string) ([]model.AgentStatus, error) {
	var _taskid string
	if taskid == nil {
		_taskid = uuid.NewString()
		taskid = &_taskid
	} else {
		_taskid = *taskid
	}

	var message = corepb.TaskDetail{
		TaskId: _taskid,
		Task: &corepb.TaskDetail_GetAgentStatusRequest{
			GetAgentStatusRequest: &gateway.GetAgentStatusRequest{
				AgentId: ids,
			},
		},
	}

	pbData, err := proto.Marshal(&message)
	if err != nil {
		return nil, exception.ErrNewException(err, exceptionCode.ErrorEncodeProtoMessage, exception.WithMsg("Failed to marshal agent information by proto"))
	}
	responseTopic, _, err := _taskHelper(
		ctx,
		a.txnHelper,
		a.store.Mq,
		a.db.AddAsyncJob,
		taskid,
		user.RoleGateway,
		gatewayid,
		string(taskTypes.TargetGateway),
		fmt.Sprintf("/task/%s", taskTypes.TargetGateway),
		pbData,
	)
	if err != nil {
		return nil, exception.ErrNewException(err, exceptionCode.ErrorPushTaskFailed, exception.WithMsg("Failed to add async task"))
	}

	var response corepb.QueryTaskResultResponse
	respDataCh, _ := a.store.Mq.SubscribeBytes(responseTopic)
	defer func(Mq mq.Mq, topic string, sub <-chan any) {
		err := Mq.Unsubscribe(topic, sub)
		if err != nil {
			logger.Errorf("failed to unsubscribe from message queue: %v", err)
		}
	}(a.store.Mq, responseTopic, respDataCh)
	select {
	case <-time.After(10 * time.Second):
		return nil, exception.ErrNewException(err, exceptionCode.ErrorPushTaskFailed, exception.WithMsg("Failed to add async task"), exception.WithMsg("timeout"))
	case m := <-respDataCh:
		err = proto.Unmarshal(m.([]byte), &response)
		if err != nil {
			return nil, exception.ErrNewException(err, exceptionCode.ErrorDecodeProtoMessage, exception.WithMsg("Failed to unmarshal agent information by proto"))
		}
	}

	var status status.Status
	result, ok := response.Result.(*corepb.QueryTaskResultResponse_Finish)
	if !ok || (ok && !result.Finish.MessageIs(&status)) {
		return nil, exception.ErrNewException(err, exceptionCode.ErrorDecodeProtoMessage, exception.WithMsg("Failed to unmarshal agent information by proto"))
	}
	err = result.Finish.UnmarshalTo(&status)
	if err != nil {
		return nil, exception.ErrNewException(err, exceptionCode.ErrorDecodeProtoMessage, exception.WithMsg("Failed to unmarshal agent information by proto"))
	}

	var retVal = make([]model.AgentStatus, len(status.GetDetails()))
	for i, detail := range status.Details {
		var s wrapperspb.StringValue
		if err := detail.UnmarshalTo(&s); err != nil {
			logger.Errorf("Failed to unmarshal agent information by proto")
			continue
		}
		retVal[i] = model.AgentStatus(s.GetValue())
	}
	return retVal, nil
}
