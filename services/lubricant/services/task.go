package services

import (
	"context"
	"fmt"

	"github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/pkg/types/user"
	"google.golang.org/protobuf/proto"
)

func (a *AgentService) PushTaskAgent(ctx context.Context, taskid *string, userID, gatewayID, agentID string, bin []byte) (string, string, error) {
	return _taskHelper(&taskArgs{
		ctx:           ctx,
		txnHelper:     a.txnHelper,
		storeMq:       a.store.Mq,
		dbAddAsyncJob: a.db.AddAsyncJob,
		taskID:        taskid,
		executorType:  user.RoleAgent,
		executorID:    gatewayID,
		userID:        userID,
		targetName:    task.TargetGateway,
		topicPrefix:   fmt.Sprintf("/task/%s", userID),
		bin:           bin,
	})
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
	return _taskHelper(&taskArgs{
		ctx:           ctx,
		txnHelper:     s.txnHelper,
		storeMq:       s.store.Mq,
		dbAddAsyncJob: s.db.AddAsyncJob,
		taskID:        taskid,
		executorType:  user.RoleGateway,
		executorID:    gatewayID,
		userID:        userID,
		targetName:    task.TargetGateway,
		topicPrefix:   fmt.Sprintf("/task/%s", userID),
		bin:           bin,
	})
}
func (s *GatewayService) PushTaskPb(ctx context.Context, taskid *string, userID, gatewayID string, pb proto.Message) (string, string, error) {
	bin, err := proto.Marshal(pb)
	if err != nil {
		return "", "", err
	}
	return s.PushTask(ctx, taskid, gatewayID, userID, bin)
}
