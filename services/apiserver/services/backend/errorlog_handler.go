package backend

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
)

var errCh = make(chan *model.ErrorLogs, 3)

type ErrLogCollect struct {
	dataCli *datastore.DataStore
}

func (e *ErrLogCollect) handler() {
	for err := range errCh {
		err := e.dataCli.ICoreDb.SaveErrorLog(context.Background(), err)
		if err != nil {
			logger.Errorf("failed to save error log: %v", err)
		}
		// todo: need to report error to user
		logger.Errorf("%v", e)
	}
}
