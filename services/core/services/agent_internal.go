package services

import (
	"github.com/AEnjoy/IoT-lubricant/services/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/services/core/models"
)

type AgentService struct {
	db    models.ICoreDb
	store *datastore.DataStore
}
