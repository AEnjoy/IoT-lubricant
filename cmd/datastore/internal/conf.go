package internal

import (
	"strconv"
	"strings"

	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
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
func GetPort() string {
	return strconv.Itoa(config.SystemConfig.GrpcPort)
}
func GetTLS() *crypto.Tls {
	return &config.SystemConfig.Tls
}
