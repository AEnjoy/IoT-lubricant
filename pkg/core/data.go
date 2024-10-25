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
		userID, _ := dataCli.GatewayIDGetUserID(context.Background(), info.GatewayId)
		compressor, err := compress.NewCompressor(info.Algorithm)
		if err != nil {
			errCh <- &ErrLogInfo{User: userID, Agent: data.GetAgentID(), Message: err}
			return
		}
		for i, in := range data.GetData() {
			decompress, err := compressor.Decompress(in)
			if err != nil {
				errCh <- &ErrLogInfo{User: userID, Agent: data.GetAgentID(), Message: err}
				return
			}

			out, err := cleaner.Run(decompress)
			if err != nil {
				errCh <- &ErrLogInfo{User: userID, Agent: data.GetAgentID(), Message: err}
				return
			}

			out, _ = compressor.Compress(out)
			data.Data[i] = out
		}
	}

	s := data.String()
	ctx := context.Background()
	_ = dataCli.HSet(ctx, data.GetAgentID(), "latest", s)
	// handel cache error is not need
	err = dataCli.StoreAgentGatherData(ctx, nil, data.GetAgentID(), s)
	if err != nil {
		info, _ := dataCli.GetAgentInfo(data.GetAgentID())
		userID, _ := dataCli.GatewayIDGetUserID(context.Background(), info.GatewayId)
		errCh <- &ErrLogInfo{User: userID, Agent: data.GetAgentID(), Message: err}
	}
}
