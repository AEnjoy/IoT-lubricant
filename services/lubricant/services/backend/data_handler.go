package backend

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/utils/compress"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"google.golang.org/protobuf/proto"
)

type DataHandler struct {
	dataCli *datastore.DataStore
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
