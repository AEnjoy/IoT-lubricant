package services

import (
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/repo"
	"github.com/aenjoy/iot-lubricant/services/corepkg/syncQueue"
)

var _ IAgentService = (*AgentService)(nil)

type AgentService struct {
	db    repo.ICoreDb
	store *datastore.DataStore
	*syncQueue.SyncTaskQueue
}
