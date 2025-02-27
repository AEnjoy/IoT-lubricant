package services

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/repo"
)

type AgentService struct {
	db    repo.ICoreDb
	store *datastore.DataStore
}
