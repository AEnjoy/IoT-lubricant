package core

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/pkg/core/datastore"
	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/compress"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
)

var dataCli = func() *datastore.DataStore {
	return ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
}()

func HandelRecvData(data *core.Data) {
	cleaner, err := dataCli.GetDataCleaner(data.GetAgentID())
	if err == nil {
		info, _ := dataCli.GetAgentInfo(data.GetAgentID())
		compressor, err := compress.NewCompressor(info.Algorithm)
		if err != nil {
			errCh <- &ErrLogInfo{User: info.UserId, Agent: data.GetAgentID(), Message: err}
			return
		}
		for i, in := range data.GetData() {
			decompress, err := compressor.Decompress(in)
			if err != nil {
				errCh <- &ErrLogInfo{User: info.UserId, Agent: data.GetAgentID(), Message: err}
				return
			}

			out, err := cleaner.Run(decompress)
			if err != nil {
				errCh <- &ErrLogInfo{User: info.UserId, Agent: data.GetAgentID(), Message: err}
				return
			}

			out, _ = compressor.Compress(out)
			data.Data[i] = out
		}
	}

	s := data.String()
	_ = dataCli.HSet(context.Background(), data.GetAgentID(), "latest", s)
	_ = dataCli.StoreAgentGatherData(data.GetAgentID(), s)
}
