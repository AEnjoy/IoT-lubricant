package backend

import (
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
)

func (ErrLogCollect) Weight() uint16 {
	return ioc.BackendHandlerErrLogs
}

func (ErrLogCollect) Version() string {
	return ""
}

func (e *ErrLogCollect) Init() error {
	e.dataCli = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	go e.handler()
	return nil
}

func (d *DataHandler) Init() error {
	d.dataCli = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	go d.handler()
	return nil
}

func (DataHandler) Weight() uint16 {
	return ioc.BackendHandlerDataUpload
}

func (d DataHandler) Version() string {
	return ""
}
