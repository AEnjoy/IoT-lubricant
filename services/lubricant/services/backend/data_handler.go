package backend

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/utils/compress"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"

	"google.golang.org/protobuf/proto"
)

var _ ioc.Object = (*DataHandler)(nil)

type DataHandler struct {
	dataCli *datastore.DataStore
}

func (d *DataHandler) Init() error {
	d.dataCli = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	go d.handler()
	return nil
}
func (d *DataHandler) handler() {
	sub, err := d.dataCli.Mq.SubscribeBytes("/handler/data")
	if err != nil {
		panic(err)
	}
	for payload := range sub {
		go d._dataStoreExecute(payload)
	}
}
func (d *DataHandler) _dataStoreExecute(payload any) {
	var data corepb.Data
	err := proto.Unmarshal(payload.([]byte), &data)
	if err != nil {
		logger.Errorf("failed to unmarshal data: %v", err)
		return
	}
	cleaner, err := d.dataCli.GetDataCleaner(data.GetAgentID())
	if err == nil {
		info, _ := d.dataCli.GetAgentInfo(data.GetAgentID())
		compressor, err := compress.NewCompressor(info.Algorithm)
		if err != nil {
			errCh <- &model.ErrorLogs{Component: "core", Module: "compressor", Message: err.Error()}
			return
		}
		for i, in := range data.GetData() {
			decompress, err := compressor.Decompress(in)
			if err != nil {
				errCh <- &model.ErrorLogs{Component: "core", Module: "compressor", Message: err.Error()}
				return
			}

			out, err := cleaner.Run(decompress)
			if err != nil {
				errCh <- &model.ErrorLogs{Component: "core", Module: "data-cleaner", Message: err.Error()}
				return
			}

			out, _ = compressor.Compress(out)
			data.Data[i] = out
		}
	}

	s := data.String()
	ctx := context.Background()
	_ = d.dataCli.HSet(ctx, data.GetAgentID(), "latest", s)
	// handel cache error is not need
	err = d.dataCli.StoreAgentGatherData(ctx, nil, data.GetAgentID(), s)
	if err != nil {
		errCh <- &model.ErrorLogs{Component: "core", Module: "datastore", Message: err.Error()}
	}

}
func (DataHandler) Weight() uint16 {
	return ioc.BackendHandlerDataUpload
}

func (d DataHandler) Version() string {
	return ""
}
