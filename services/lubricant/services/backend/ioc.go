package backend

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
)

var (
	_ ioc.Object = (*ReportHandler)(nil)
	_ ioc.Object = (*GatewayGuard)(nil)
	_ ioc.Object = (*ErrLogCollect)(nil)
	_ ioc.Object = (*DataHandler)(nil)
)

func (ReportHandler) Weight() uint16 {
	return ioc.BackendHandlerReport
}

func (r *ReportHandler) Init() error {
	r.dataCli = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	go r.handler()
	return nil
}
func (ReportHandler) Version() string {
	return ""
}

func (GatewayGuard) Weight() uint16 {
	return ioc.GatewayStatusGuard
}

func (GatewayGuard) Version() string {
	return ""
}

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
