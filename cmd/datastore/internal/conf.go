package internal

import (
	"strings"

	"github.com/aenjoy/iot-lubricant/services/corepkg/config"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
)

func GetDataStore() *datastore.DataStore {
	return ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
}
func GetEtcdEndpoints() []string {
	return strings.Split(config.GetConfig().EtcdEndpoints, ",")
}
func GetInternalWorkThreadNumber() int {
	return config.GetConfig().InternalWorkThreadNumber
}
