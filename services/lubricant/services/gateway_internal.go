package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/model/request"
	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	proxypb "github.com/aenjoy/iot-lubricant/protobuf/gateway"
	"github.com/aenjoy/iot-lubricant/services/lubricant/datastore"
	"github.com/aenjoy/iot-lubricant/services/lubricant/repo"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/status"
)

var _ IGatewayService = (*GatewayService)(nil)

type GatewayService struct {
	db    repo.ICoreDb
	store *datastore.DataStore
	*SyncTaskQueue
}

func (s *GatewayService) EditGateway(ctx context.Context, gatewayid, description string, tls *crypto.Tls) error {
	txn, errorCh, commit := s.txnHelper()
	defer commit()

	var info model.Gateway
	if tls != nil {
		marshalString, err := sonic.MarshalString(tls)
		if err != nil {
			err = exception.ErrNewException(err,
				exceptionCode.ErrorEncodeJSON,
				exception.WithMsg("Failed to marshall json result"),
			)
			errorCh.Report(err, exceptionCode.ErrorEncodeJSON, "Failed to marshall json result", true)
			return err
		}
		info.TlsConfig = marshalString
	}

	info.Description = description
	info.GatewayID = gatewayid
	err := s.db.UpdateGateway(ctx, txn, info)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbUpdateGatewayInfoFailed,
			exception.WithMsg("Failed to update gateway information to database"),
		)
		return err
	}
	return nil
}

func (s *GatewayService) DescriptionGateway(ctx context.Context, userid, gatewayid string) (*response.DescriptionGatewayResponse, error) {
	var retVal response.DescriptionGatewayResponse
	info, err := s.db.GetGatewayInfo(ctx, gatewayid)
	if err != nil {
		return &retVal, exception.ErrNewException(err, exceptionCode.DbGetGatewayFailed, exception.WithMsg("Failed to get gateway information from database"))
	}
	retVal.Gateway = info

	agents, err := s.db.GetAgentList(ctx, userid, gatewayid)
	if err != nil {
		return &retVal, exception.ErrNewException(err, exceptionCode.GetAgentFailed, exception.WithMsg("Failed to get gateway information from database"))
	}
	retVal.Agents = agents
	return &retVal, nil
}

func (s *GatewayService) DescriptionHost(ctx context.Context, hostid string) (*response.DescriptionHostResponse, error) {
	var retVal response.DescriptionHostResponse
	info, err := s.db.GetGatewayHostInfo(ctx, hostid)
	if err != nil {
		return &retVal, exception.ErrNewException(err, exceptionCode.DbGetGatewayFailed, exception.WithMsg("Failed to get gateway host information from database"))
	}
	retVal.Host = info
	gateways, err := s.db.GetAllGatewayByUserID(ctx, info.UserID)
	if err != nil {
		return &retVal, exception.ErrNewException(err, exceptionCode.DbGetGatewayFailed, exception.WithMsg("Failed to get gateway information from database"))
	}
	retVal.GatewayList = gateways
	return &retVal, nil
}

func (s *GatewayService) AddHostInternal(ctx context.Context, info *model.GatewayHost) error {
	txn, _, commit := s.txnHelper()
	defer commit()

	err := s.db.AddGatewayHostInfo(ctx, txn, info)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbAddGatewayFailed,
			exception.WithMsg("Failed to add gateway information to database"),
		)
		return err
	}
	return nil
}
func (s *GatewayService) AddGatewayInternal(ctx context.Context, userid, gatewayid, description string, tls *crypto.Tls) error {
	txn, errorCh, commit := s.txnHelper()
	defer commit()

	err := s.db.AddGateway(ctx, txn, userid, model.Gateway{
		GatewayID:   gatewayid,
		Description: description,
		TlsConfig: func() string {
			if tls == nil {
				return ""
			}
			marshalString, err := sonic.MarshalString(tls)
			if err != nil {
				return ""
			}
			return marshalString
		}(),
	})
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbAddGatewayFailed,
			exception.WithMsg("Failed to add gateway information to database"),
		)
		errorCh.Report(err, exceptionCode.DbAddGatewayFailed, "Failed to add gateway information to database", true)
		return err
	}
	return nil
}
func (s *GatewayService) RemoveGatewayInternal(ctx context.Context, gatewayid string) error {
	txn, _, commit := s.txnHelper()
	defer commit()

	taskMq := s.store.Mq
	topic := fmt.Sprintf("/monitor/%s/%s/unregister", taskTypes.TargetGateway, gatewayid)
	err := taskMq.PublishBytes(topic, []byte("unregister gateway"))
	if err != nil {
		return exception.ErrNewException(err,
			exceptionCode.MqPublishFailed,
			exception.WithMsg("Failed to send gateway unregister signal"),
		)
	}
	t, err := taskMq.Subscribe(fmt.Sprintf("%s/response", topic))
	if err != nil {
		return exception.ErrNewException(err,
			exceptionCode.MqSubscribeFailed,
			exception.WithMsg("Failed to set gateway unregister signal"),
		)
	}

	select {
	case <-ctx.Done():
		return exception.ErrNewException(ctx.Err(),
			exceptionCode.DeadLine,
			exception.WithMsg("user request cancel"),
			exception.WithMsg("database not changed"))
	case m := <-t:
		if string(m.([]byte)) != "ok" {
			return exception.ErrNewException(err,
				exceptionCode.RemoveGatewayFailed,
				exception.WithMsg("gateway monitor failed to unregister this gateway"),
				exception.WithMsg(fmt.Sprintf("GatewayID: %s", gatewayid)),
				exception.WithMsg(fmt.Sprintf("Message: %s", m)),
			)
		}
	case <-time.After(10 * time.Second):
		return exception.ErrNewException(err,
			exceptionCode.RemoveGatewayFailed,
			exception.WithMsg("gateway monitor failed to unregister this gateway"),
			exception.WithMsg(fmt.Sprintf("GatewayID: %s", gatewayid)),
			exception.WithMsg("timeout"),
		)
	}
	err = s.db.DeleteGateway(ctx, txn, gatewayid)
	if err != nil {
		return exception.ErrNewException(err,
			exceptionCode.RemoveGatewayFailed,
			exception.WithMsg("Failed to delete gateway information from database"),
		)
	}
	return nil
}
func (s *GatewayService) RemoveGatewayHostInternal(ctx context.Context, hostid string) error {
	panic("implement me")
}

func (s *GatewayService) AddAgentInternal(ctx context.Context, taskid *string, userid, gatewayid string,
	req *request.AddAgentRequest, openapidoc, enableFile []byte) (string, error) {
	txn, errorCh, commit := s.txnHelper()

	var conf = model.CreateAgentConf{AgentContainerInfo: req.AgentContainerInfo, DriverContainerInfo: req.DriverContainerInfo}
	confData, err := sonic.Marshal(&conf)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.ErrorEncodeJSON,
			exception.WithMsg("Failed to marshal agent information"),
		)
		return "", err
	}

	var (
		td      corepb.TaskDetail
		pb      proxypb.CreateAgentRequest
		_true   = true
		agentID = uuid.NewString()
	)

	pb.Info = &agentpb.AgentInfo{
		AgentID:     agentID,
		GatewayID:   &gatewayid,
		Description: &req.Description,
		GatherCycle: &req.GatherCycle,
		Algorithm:   &req.DataCompressAlgorithm,
		DataSource: &agentpb.OpenapiDoc{
			OriginalFile: openapidoc,
			EnableFile:   enableFile,
		},
		Stream: &req.EnableStreamAbility,

		ReportCycle: &req.ReportCycle,
		Address:     &req.Address,
	}
	pb.Conf = confData

	td.TaskId = *taskid
	td.IsSynchronousTask = &_true
	td.Task = &corepb.TaskDetail_CreateAgentRequest{
		CreateAgentRequest: &pb,
	}

	err = s.db.AddAgent(ctx, txn, gatewayid, model.Agent{
		AgentId:     agentID,
		GatewayId:   gatewayid,
		Description: req.Description,
		Algorithm:   req.DataCompressAlgorithm,
		GatherCycle: int(req.GatherCycle),
	})
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.AddAgentFailed,
			exception.WithMsg("Failed to add agent information to database"),
		)
		return "", err
	}

	_, _, err = s.PushTaskPb(ctx, taskid, userid, gatewayid, &td)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.ErrorPushTaskFailed,
			exception.WithMsg("Failed to push agent information to gateway"),
		)
		errorCh.Report(err, exceptionCode.ErrorPushTaskFailed, "Failed to push agent information to gateway", false)
		return "", err
	}

	go func() {
		defer commit()
		resp, err := s.SyncTaskQueue.WaitTask(*taskid, 10*time.Second)
		if err != nil {
			errorCh.Report(err, exceptionCode.AddAgentFailed, "timeout", false)
			return
		}
		if resp.GetFinish() != nil {
			var s status.Status
			if err := resp.GetFinish().UnmarshalTo(&s); err != nil {
				errorCh.Report(err, exceptionCode.AddAgentFailed, "Failed to add agent information to database", false)
				return
			}
			if len(s.Details) == 0 {
				errorCh.Report(err, exceptionCode.AddAgentFailed, "Failed to add agent information to database", false)
				return
			}
		}
	}()
	return agentID, nil
}
