package dataapi

import (
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
	"github.com/aenjoy/iot-lubricant/services/corepkg/config"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
)

var _ ioc.Object = (*DataStoreApiService)(nil)

type DataStoreApiService struct {
	svcpb.DataStoreServiceClient
	CloseCall func() error
}

func (d *DataStoreApiService) Init() (err error) {
	c := config.GetConfig()
	if c.SvcDataStoreMode == "rpc" {
		logger.Info("DataStoreApiService.Init() rpc mode")
		d.DataStoreServiceClient, d.CloseCall, err =
			svcpb.NewDataStoreServiceClientWithHost(
				c.SvcDataStoreEndpoint,
				c.SvcDataStoreTls)
	}

	return
}

func (DataStoreApiService) Weight() uint16 {
	return ioc.SvcDataStoreApiService
}

func (DataStoreApiService) Version() string {
	return ""
}
