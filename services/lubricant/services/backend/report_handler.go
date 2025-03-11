package backend

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
	"google.golang.org/protobuf/proto"
)

var _ ioc.Object = (*ReportHandler)(nil)

type ReportHandler struct {
	dataCli *datastore.DataStore
}

func (r *ReportHandler) Init() error {
	r.dataCli = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE_STORE).(*datastore.DataStore)
	go r.handler()
	return nil
}
func (r *ReportHandler) handler() {
	sub, err := r.dataCli.Mq.SubscribeBytes("/handler/report")
	if err != nil {
		panic(err)
	}
	for payload := range sub {
		go r._reportPayload(payload)
	}
}
func (r *ReportHandler) _reportPayload(payload any) {
	var req corepb.ReportRequest
	err := proto.Unmarshal(payload.([]byte), &req)
	if err != nil {
		logger.Errorf("failed to unmarshal report: %v", err)
		return
	}
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
		txn := r.dataCli.Begin()
		err := r.dataCli.UpdateAgentStatus(context.Background(), txn, req.GetAgentId(), req.GetAgentStatus().GetReq().GetMessage())
		if err != nil {
			logger.Errorf("failed to update agent status: %v", err)
		}
		r.dataCli.Commit(txn)
	}
}
func (ReportHandler) Weight() uint16 {
	return ioc.BackendHandlerReport
}

func (ReportHandler) Version() string {
	return ""
}
