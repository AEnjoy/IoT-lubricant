package services

import (
	"context"
	"fmt"

	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/pkg/types/user"
	"google.golang.org/protobuf/proto"
)

func (a *AgentService) PushTaskAgent(ctx context.Context, taskid *string, userID, gatewayID, agentID string, bin []byte) (string, string, error) {
	return _taskHelper(
		ctx,
		a.txnHelper,
		a.store.Mq,
		a.db.AddAsyncJob,
		taskid,
		user.RoleAgent,
		gatewayID, //agentID,
		userID,
		string(taskTypes.TargetAgent),
		fmt.Sprintf("/task/%s/%s", userID, taskTypes.TargetGateway), //fmt.Sprintf("/task/%s/%s/%s", taskTypes.TargetGateway, gatewayID, taskTypes.TargetAgent),
		bin,
	)
}

// PushTaskAgentPb :
// pb is core.TaskDetail
func (a *AgentService) PushTaskAgentPb(ctx context.Context, taskid *string, userID, gatewayID, agentID string, pb proto.Message) (string, string, error) {
	bin, err := proto.Marshal(pb)
	if err != nil {
		return "", "", err
	}
	return a.PushTaskAgent(ctx, taskid, userID, gatewayID, agentID, bin)
}

func (s *GatewayService) PushTask(ctx context.Context, taskid *string, gatewayID, userID string, bin []byte) (string, string, error) {
	return _taskHelper(
		ctx,
		s.txnHelper,
		s.store.Mq,
		s.db.AddAsyncJob,
		taskid,
		user.RoleGateway,
		gatewayID,
		userID,
		string(taskTypes.TargetGateway),
		fmt.Sprintf("/task/%s/%s", userID, taskTypes.TargetGateway),
		bin,
	)
}
func (s *GatewayService) PushTaskPb(ctx context.Context, taskid *string, userID, gatewayID string, pb proto.Message) (string, string, error) {
	bin, err := proto.Marshal(pb)
	if err != nil {
		return "", "", err
	}
	return s.PushTask(ctx, taskid, gatewayID, userID, bin)
}
