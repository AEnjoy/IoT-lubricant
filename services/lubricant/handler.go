// Package lubricant handler.go
// Handler.go is a logical collection used for asynchronous data processing
// The processing content includes GRPC data stream decoupling, asynchronous task decoupling, etc
package lubricant

import (
	"context"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/utils/compress"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
)

var dataCli = func() *datastore.DataStore {
	return ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
}

func HandelRecvData(data *corepb.Data) {
	// todo: use mq to handel data
	dataCli := dataCli()
	cleaner, err := dataCli.GetDataCleaner(data.GetAgentID())
	if err == nil {
		info, _ := dataCli.GetAgentInfo(data.GetAgentID())
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
	_ = dataCli.HSet(ctx, data.GetAgentID(), "latest", s)
	// handel cache error is not need
	err = dataCli.StoreAgentGatherData(ctx, nil, data.GetAgentID(), s)
	if err != nil {
		errCh <- &model.ErrorLogs{Component: "core", Module: "datastore", Message: err.Error()}
	}
}

func HandelReport(req *corepb.ReportRequest) {
	// todo: use mq to handel data
	dataStore := dataCli()
	switch data := req.GetReq().(type) {
	case *corepb.ReportRequest_Error:
		e := data.Error.GetErrorMessage()
		code := e.GetCode()

		errCh <- &model.ErrorLogs{
			Component: e.GetModule(),
			Code: func() int32 {
				if code == nil {
					return 0
				}
				return code.GetCode()
			}(),
			Type: e.GetErrorType(),
			Message: func() string {
				if code == nil {
					return ""
				}
				return code.GetMessage()
			}(),
			Stack: e.GetStack(),
		}
	case *corepb.ReportRequest_AgentStatus:
		txn := dataStore.Begin()
		err := dataStore.UpdateAgentStatus(context.Background(), txn, req.GetAgentId(), req.GetAgentStatus().GetReq().GetMessage())
		if err != nil {
			logger.Errorf("failed to update agent status: %v", err)
		}
		dataStore.Commit(txn)
	}
}

var errCh = make(chan *model.ErrorLogs, 3)

func handleErrLog() {
	time.Sleep(3 * time.Second)
	dataStore := dataCli()
	for e := range errCh {
		err := dataStore.ICoreDb.SaveErrorLog(context.Background(), e)
		if err != nil {
			logger.Errorf("failed to save error log: %v", err)
		}
		// todo: need to report error to user
		logger.Errorf("%v", e)
	}
}
func init() {
	go gatewayStatusGuard()
	go handleErrLog()
}
