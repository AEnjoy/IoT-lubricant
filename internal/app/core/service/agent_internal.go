package service

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
)

type AgentService struct {
	db    repo.ICoreDb
	store *datastore.DataStore
}
