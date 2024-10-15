package app

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/cmd/core/app/datastore"
	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
)

var dataCli = func() *datastore.DataStore {
	return ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
}()

func HandelRecvData(data *core.Data) {
	s := data.String()
	_ = dataCli.HSet(context.Background(), data.GetAgentID(), "latest", s)
	_ = dataCli.StoreAgentGatherData(data.GetAgentID(), s)
}
