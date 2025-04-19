package backend

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	errChHandel "github.com/aenjoy/iot-lubricant/pkg/error"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/aenjoy/iot-lubricant/protobuf/svc"
	"github.com/aenjoy/iot-lubricant/services/apiserver/services"
	"github.com/aenjoy/iot-lubricant/services/corepkg/datastore"

	"google.golang.org/protobuf/proto"
)

type ReportHandler struct {
	dataCli *datastore.DataStore
	*services.SyncTaskQueue
	//*services.AgentService
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
	ctx := context.Background()
	txn := r.dataCli.ICoreDb.Begin()
	errorCh := errChHandel.NewErrorChan()
	defer errChHandel.HandleErrorCh(errorCh).
		ErrorWillDo(func(error) {
			r.dataCli.ICoreDb.Rollback(txn)
		}).
		SuccessWillDo(func() {
			r.dataCli.ICoreDb.Commit(txn)
		}).
		Do()

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
		err := r.dataCli.UpdateAgentStatus(ctx, txn, req.GetAgentId(), req.GetAgentStatus().GetReq().GetMessage())
		if err != nil {
			errorCh.Report(err, exceptionCode.UpdateAgentStatusFailed, "database error", true)
		}
	case *corepb.ReportRequest_TaskResult:
		msg := data.TaskResult.GetMsg()
		taskid := msg.GetTaskId()
		if err := r.SyncTaskQueue.FinshTask(taskid, msg); err != nil {
			logger.Errorf("Failed to send task finsh result to message middleware: %v", err)
		}

		switch result := msg.GetResult().(type) {
		case *corepb.QueryTaskResultResponse_Finish:
			if err := r.dataCli.SetAsyncJobStatus(ctx, txn, taskid, "completed", result.Finish.String()); err != nil {
				errorCh.Report(err, exceptionCode.UpdateTaskStatusFailed, "database error", true)
			}
		case *corepb.QueryTaskResultResponse_Failed:
			if err := r.dataCli.SetAsyncJobStatus(ctx, txn, taskid, "failed", result.Failed.String()); err != nil {
				errorCh.Report(err, exceptionCode.UpdateTaskStatusFailed, "database error", true)
			}
		case *corepb.QueryTaskResultResponse_Pending:
			if err := r.dataCli.SetAsyncJobStatus(ctx, txn, taskid, "pending", result.Pending.String()); err != nil {
				errorCh.Report(err, exceptionCode.UpdateTaskStatusFailed, "database error", true)
			}
		default:
			if err := r.dataCli.SetAsyncJobStatus(ctx, txn, taskid, "started", ""); err != nil {
				errorCh.Report(err, exceptionCode.UpdateTaskStatusFailed, "database error", true)
			}
		}
	case *corepb.ReportRequest_ReportLog:
		for _, logs := range data.ReportLog.GetLogs() {
			go func(logs *svc.Logs) {
				data, err := proto.Marshal(logs)
				if err != nil {
					logger.Errorf("failed to marshal logs: %v", err)
				}
				err = r.dataCli.Mq.PublishBytes(constant.MESSAGE_SVC_LOGGER, data)
				if err != nil {
					logger.Errorf("failed to publish logs: %v", err)
				}
			}(logs)
		}
	default:
	}
}
