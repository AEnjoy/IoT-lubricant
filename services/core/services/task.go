package services

import (
	"context"
	"fmt"

	taskTypes "github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/user"
)

func (a *AgentService) PushTaskAgent(ctx context.Context, taskid *string, gatewayID, agentID string, bin []byte) (string, string, error) {
	return _taskHelper(
		ctx,
		a.txnHelper,
		a.store.Mq,
		a.db.AddAsyncJob,
		taskid,
		user.RoleAgent,
		agentID,
		string(taskTypes.TargetAgent),
		fmt.Sprintf("/task/%s/%s/%s", taskTypes.TargetGateway, gatewayID, taskTypes.TargetAgent),
		bin,
	)
}

func (s *GatewayService) PushTask(ctx context.Context, taskid *string, gatewayID string, bin []byte) (string, string, error) {
	return _taskHelper(
		ctx,
		s.txnHelper,
		s.store.Mq,
		s.db.AddAsyncJob,
		taskid,
		user.RoleGateway,
		gatewayID,
		string(taskTypes.TargetGateway),
		fmt.Sprintf("/task/%s", taskTypes.TargetGateway),
		bin,
	)
}
