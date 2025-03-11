// Package lubricant handler.go
// Handler.go is a logical collection used for asynchronous data processing
// The processing content includes GRPC data stream decoupling, asynchronous task decoupling, etc
package lubricant

import (
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
	"google.golang.org/protobuf/proto"
)

func HandelRecvData(data *corepb.Data) {
	dataCli := ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	marshal, err := proto.Marshal(data)
	if err != nil {
		logger.Errorf("failed to marshal data: %v", err)
	}

	err = dataCli.Mq.PublishBytes("/handler/data", marshal)
	if err != nil {
		logger.Errorf("failed to publish data: %v", err)
	}
}
