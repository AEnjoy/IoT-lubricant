package services

import (
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/repo"
)

var _ IAgentService = (*AgentService)(nil)

type AgentService struct {
	db    repo.ICoreDb
	store *datastore.DataStore
	*SyncTaskQueue
}
