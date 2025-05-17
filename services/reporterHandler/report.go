package reporterHandler

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	errChHandel "github.com/aenjoy/iot-lubricant/pkg/error"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/aenjoy/iot-lubricant/protobuf/svc"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"google.golang.org/protobuf/proto"
)

func (a app) reporterPayload(payload any) {
	var req corepb.ReportRequest
	err := proto.Unmarshal(payload.([]byte), &req)
	if err != nil {
		logg.L.Errorf("failed to unmarshal report: %v", err)
		return
	}
	ctx := context.Background()
	txn := a.ICoreDb.Begin()
	errorCh := errChHandel.NewErrorChan()
	defer errChHandel.HandleErrorCh(errorCh).
		ErrorWillDo(func(error) {
			a.ICoreDb.Rollback(txn)
		}).
		SuccessWillDo(func() {
			a.ICoreDb.Commit(txn)
		}).
		Do()

	switch data := req.GetReq().(type) {
	case *corepb.ReportRequest_Error:
		e := data.Error.GetErrorMessage()
		code := e.GetCode()
		logg.L.
			WithAction("DevicesUploadError").
			WithNotPrintToStdout().
			WithExceptionCode(exceptionCode.ResCode(code.GetCode())).
			WithOperationType(svc.Operation(e.GetErrorType())).
			WithMetaData(e).
			Error(code.GetMessage())

	case *corepb.ReportRequest_AgentStatus:
		err := a.UpdateAgentStatus(ctx, txn, req.GetAgentId(), req.GetAgentStatus().GetReq().GetMessage())
		if err != nil {
			errorCh.Report(err, exceptionCode.UpdateAgentStatusFailed, "database error", true)
		}
	case *corepb.ReportRequest_TaskResult:
		msg := data.TaskResult.GetMsg()
		taskid := msg.GetTaskId()
		if err := a.SyncTaskQueue.FinshTask(taskid, msg); err != nil {
			logger.Errorf("Failed to send task finsh result to message middleware: %v", err)
		}

		switch result := msg.GetResult().(type) {
		case *corepb.QueryTaskResultResponse_Finish:
			if err := a.SetAsyncJobStatus(ctx, txn, taskid, "completed", result.Finish.String()); err != nil {
				errorCh.Report(err, exceptionCode.UpdateTaskStatusFailed, "database error", true)
			}
		case *corepb.QueryTaskResultResponse_Failed:
			if err := a.SetAsyncJobStatus(ctx, txn, taskid, "failed", result.Failed.String()); err != nil {
				errorCh.Report(err, exceptionCode.UpdateTaskStatusFailed, "database error", true)
			}
		case *corepb.QueryTaskResultResponse_Pending:
			if err := a.SetAsyncJobStatus(ctx, txn, taskid, "pending", result.Pending.String()); err != nil {
				errorCh.Report(err, exceptionCode.UpdateTaskStatusFailed, "database error", true)
			}
		default:
			if err := a.SetAsyncJobStatus(ctx, txn, taskid, "started", ""); err != nil {
				errorCh.Report(err, exceptionCode.UpdateTaskStatusFailed, "database error", true)
			}
		}
	case *corepb.ReportRequest_ReportLog:
		for _, logs := range data.ReportLog.GetLogs() {
			go func(logs *svc.Logs) {
				data, err := proto.Marshal(logs)
				if err != nil {
					logg.L.Errorf("failed to marshal logs: %v", err)
				}
				err = a.DataStore.V2mq.QueuePublish(constant.MESSAGE_SVC_LOGGER, data)
				if err != nil {
					logg.L.Errorf("failed to publish logs: %v", err)
				}
			}(logs)
		}
	default:
	}
}
