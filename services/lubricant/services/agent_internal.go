package services

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/repo"
)

var _ IAgentService = (*AgentService)(nil)

type AgentService struct {
	db    repo.ICoreDb
	store *datastore.DataStore
	*SyncTaskQueue
}
