package internal

import (
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
	"github.com/aenjoy/iot-lubricant/services/corepkg/syncQueue"
)

func GetDataStore() *datastore.DataStore {
	return ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
}
func GetSyncTaskQueue() *syncQueue.SyncTaskQueue {
	return ioc.Controller.Get(ioc.APP_NAME_CORE_Internal_SyncTask_SERVICE).(*syncQueue.SyncTaskQueue)
}
