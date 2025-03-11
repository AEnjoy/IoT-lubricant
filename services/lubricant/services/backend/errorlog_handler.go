package backend

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
)

var _ ioc.Object = (*ErrLogCollect)(nil)
var errCh = make(chan *model.ErrorLogs, 3)

type ErrLogCollect struct {
	dataCli *datastore.DataStore
}

func (e *ErrLogCollect) Init() error {
	e.dataCli = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	go e.handler()
	return nil
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

func (ErrLogCollect) Weight() uint16 {
	return ioc.BackendHandlerErrLogs
}

func (ErrLogCollect) Version() string {
	return ""
}
