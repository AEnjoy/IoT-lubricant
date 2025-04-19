// Package lubricant handler.go
// Handler.go is a logical collection used for asynchronous data processing
// The processing content includes GRPC data stream decoupling, asynchronous task decoupling, etc
package lubricant

import (
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
	"google.golang.org/protobuf/proto"
)

func HandelRecvData(data *corepb.Data) {
	dataCli := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	marshal, err := proto.Marshal(data)
	if err != nil {
		logg.L.Errorf("failed to marshal data: %v", err)
	}

	err = dataCli.Mq.PublishBytes("/handler/data", marshal)
	if err != nil {
		logg.L.Errorf("failed to publish data: %v", err)
	}
}
