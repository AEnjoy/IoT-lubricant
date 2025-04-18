package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/ssh"
	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/status"
	"gopkg.in/yaml.v3"
)

func (s *GatewayService) AddHost(ctx context.Context, info *model.GatewayHost) error {
	txn, errorCh, commit := s.txnHelper()
	defer commit()

	err := s.db.AddGatewayHostInfo(ctx, txn, info)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbAddGatewayFailed,
			exception.WithMsg("Failed to add gateway information to database"),
		)
		return err
	}

	err = s.checkSSHLinker(info)
	if err != nil {
		errorCh.Report(err, exceptionCode.LinkToGatewayFailed, "LinkToTargetHostError:", true)
		return err
	}
	return nil
}
func (s *GatewayService) EditHost(ctx context.Context, hostid string, info *model.GatewayHost) error {
	txn, errorCh, commit := s.txnHelper()
	defer commit()

	hostInfo, err := s.db.GetGatewayHostInfo(ctx, hostid)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbGetGatewayFailed,
			exception.WithMsg("Failed to get gateway information from database"),
			exception.WithMsg("cannot compare gateway information"),
		)
		errorCh.Report(err, exceptionCode.DbGetGatewayFailed, "Failed to get gateway information from database", true)
		return err
	}
	if info.UserName != "" {
		hostInfo.UserName = info.UserName
	}
	if info.Host != "" {
		hostInfo.Host = info.Host
	}
	if info.PassWd != "" {
		hostInfo.PassWd = info.PassWd
	}
	if info.PrivateKey != "" {
		hostInfo.PrivateKey = info.PrivateKey
	}
	if info.Description != "" {
		hostInfo.Description = info.Description
	}
	err = s.db.UpdateGatewayHostInfo(ctx, txn, hostid, &hostInfo)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbUpdateGatewayInfoFailed,
			exception.WithMsg("Failed to update gateway information to database"),
		)
		errorCh.Report(err, exceptionCode.DbUpdateGatewayInfoFailed, "Failed to update gateway information to database", true)
		return err
	}
	return nil
}
func (s *GatewayService) GetHost(ctx context.Context, hostid string) (model.GatewayHost, error) {
	return s.db.GetGatewayHostInfo(ctx, hostid)
}
func (s *GatewayService) UserGetHosts(ctx context.Context, userid string) ([]model.GatewayHost, error) {
	return s.db.ListGatewayHostInfoByUserID(ctx, userid)
}

// DeployGatewayInstance 部署网关实例，返回gatewayID,error
func (s *GatewayService) DeployGatewayInstance(ctx context.Context,
	hostid, description string, tls *crypto.Tls) (string, error) {

	hostInfo, err := s.db.GetGatewayHostInfo(ctx, hostid)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbGetGatewayFailed,
			exception.WithMsg("Failed to get gateway information from database"),
		)
		return "", err
	}

	host, err := ssh.NewSSHClient(&hostInfo, false)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.LinkToGatewayFailed,
			exception.WithMsg("LinkToTargetHostError:"),
		)
		return "", err
	}

	gatewayID := uuid.NewString()
	serverInfo := s.getHostInfo()
	serverInfo.GatewayID = gatewayID

	err = host.DeployGateway(serverInfo)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.ErrorDeployGatewayFailed,
			exception.WithMsg("DeployGatewayFailed"),
		)
		return "", err
	}
	// todo:check gateway status

	err = s.AddGatewayInternal(ctx, hostInfo.UserID, gatewayID, description, tls)
	if err != nil {
		return "", err
	}
	return gatewayID, nil
}
func (s *GatewayService) GetRegisterStatus(_ context.Context, gatewayid string) *status.Status {
	taskMq := s.store.Mq
	topic := fmt.Sprintf("/ping/%s/%s/register", taskTypes.TargetGateway, gatewayid)
	t, err := taskMq.Subscribe(topic)
	if err != nil {
		err = fmt.Errorf("failed to get gateway register information: %v", err)
		return &status.Status{
			Code:    int32(exceptionCode.GetGatewayFailed),
			Message: err.Error(),
		}
	}
	defer func(taskMq mq.Mq, topic string, sub <-chan any) {
		err := taskMq.Unsubscribe(topic, sub)
		if err != nil {
			logger.Errorln("failed to unsubscribe from message queue: %v", err)
		}
	}(taskMq, topic, t)

	select {
	case id := <-t:
		if string(id.([]byte)) != gatewayid {
			return &status.Status{
				Code:    int32(exceptionCode.GetGatewayFailed),
				Message: "gateway is not registered(get gatewayid is not correct)",
			}
		} else {
			return &status.Status{
				Code:    int32(exceptionCode.Success),
				Message: "gateway is registered",
			}
		}
	case <-time.After(10 * time.Second):
		return &status.Status{
			Code:    int32(exceptionCode.GetGatewayFailed),
			Message: "gateway is not registered(timeout)",
		}
	}
}

func (s *GatewayService) GetStatus(_ context.Context, gatewayid string) *status.Status {
	taskMq := s.store.Mq

	// 发送随机消息，如果响应同发送一致，则认为网关status正常
	message := uuid.NewString()
	_ = taskMq.PublishBytes(fmt.Sprintf("/monitor/%s/%s/random-message", taskTypes.TargetGateway, gatewayid),
		[]byte(message))
	t, err := taskMq.Subscribe(fmt.Sprintf("/monitor/%s/%s/random-message/response", taskTypes.TargetGateway, gatewayid))
	if err != nil {
		return &status.Status{
			Code:    int32(exceptionCode.GetGatewayFailed),
			Message: fmt.Sprintf("failed to get gateway status: %v", err),
		}
	}
	defer func(taskMq mq.Mq, topic string, sub <-chan any) {
		err := taskMq.Unsubscribe(topic, sub)
		if err != nil {
			logger.Errorln("failed to unsubscribe from message queue: %v", err)
		}
	}(taskMq, fmt.Sprintf("/monitor/%s/%s/random-message/response", taskTypes.TargetGateway, gatewayid), t)

	id := <-t
	if string(id.([]byte)) == message {
		return &status.Status{
			Code:    int32(exceptionCode.Success),
			Message: "gateway is running",
		}
	} else {
		return &status.Status{
			Code:    int32(exceptionCode.GetGatewayFailed),
			Message: "gateway is not running",
		}
	}
}

func (s *GatewayService) GetErrorLogs(ctx context.Context,
	gatewayid string, from, to time.Time, limit int) ([]model.ErrorLogs, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.db.GetErrorLogs(ctx, gatewayid, from, to, limit)
}
func (s *GatewayService) DescriptionError(ctx context.Context, errorID string) (model.ErrorLogs, error) {
	return s.db.GetErrorLogByErrorID(ctx, errorID)
}

func (s *GatewayService) HostGetGatewayDeployConfig(ctx context.Context, hostid string) (*model.ServerInfo, error) {
	hostInfo, err := s.db.GetGatewayHostInfo(ctx, hostid)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbGetGatewayFailed,
			exception.WithMsg("Failed to get gateway information from database"),
		)
		return nil, err
	}

	host, err := ssh.NewSSHClient(&hostInfo, false)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.LinkToGatewayFailed,
			exception.WithMsg("LinkToTargetHostError:"),
		)
		return nil, err
	}
	defer func(host ssh.RemoteClient) {
		err := host.Close()
		if err != nil {
			logger.Errorln("failed to close ssh client: %v", err)
		}
	}(host)

	return host.GetConfig()
}
func (s *GatewayService) GatewayGetGatewayDeployConfig(ctx context.Context, gatewayid string) (*model.ServerInfo, error) {
	taskMq := s.store.Mq
	topic := fmt.Sprintf("/config/%s/%s/get/deploy-config", taskTypes.TargetGateway, gatewayid)
	_ = taskMq.Publish(topic, nil)
	t, err := taskMq.Subscribe(fmt.Sprintf("%s/response", topic))
	if err != nil {
		err = fmt.Errorf("failed to get gateway deploy config: %v", err)
		return nil, err
	}
	defer func(taskMq mq.Mq, topic string, sub <-chan any) {
		err := taskMq.Unsubscribe(topic, sub)
		if err != nil {
			logger.Errorf("failed to unsubscribe from message queue: %v", err)
		}
	}(taskMq, fmt.Sprintf("%s/response", topic), t)

	var ret model.ServerInfo
	err = sonic.Unmarshal((<-t).([]byte), &ret)
	return &ret, err
}

func (s *GatewayService) HostSetGatewayDeployConfig(ctx context.Context, hostid string, info *model.ServerInfo) error {
	hostInfo, err := s.db.GetGatewayHostInfo(ctx, hostid)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.DbGetGatewayFailed,
			exception.WithMsg("Failed to get gateway information from database"),
		)
		return err
	}

	host, err := ssh.NewSSHClient(&hostInfo, false)
	if err != nil {
		err = exception.ErrNewException(err,
			exceptionCode.LinkToGatewayFailed,
			exception.WithMsg("LinkToTargetHostError:"),
		)
		return err
	}
	defer func(host ssh.RemoteClient) {
		err := host.Close()
		if err != nil {
			logger.Errorf("failed to close ssh client: %v", err)
		}
	}(host)

	return host.UpdateConfig(info)
}

func (s *GatewayService) GatewaySetGatewayDeployConfig(ctx context.Context, gatewayid string, info *model.ServerInfo) error {
	taskMq := s.store.Mq
	topic := fmt.Sprintf("/config/%s/%s/set/deploy-config", taskTypes.TargetGateway, gatewayid)
	infoBytes, err := yaml.Marshal(info)
	if err != nil {
		return err
	}
	_ = taskMq.PublishBytes(topic, infoBytes)
	t, err := taskMq.Subscribe(fmt.Sprintf("%s/response", topic))
	if err != nil {
		err = fmt.Errorf("failed to set gateway deploy config: %v", err)
		return err
	}
	defer func(taskMq mq.Mq, topic string, sub <-chan any) {
		err := taskMq.Unsubscribe(topic, sub)
		if err != nil {
			logger.Errorf("failed to unsubscribe from message queue: %v", err)
		}
	}(taskMq, fmt.Sprintf("%s/response", topic), t)

	select {
	case resp := <-t:
		if string(resp.([]byte)) != "ok" {
			return fmt.Errorf("failed to set gateway deploy config: %s", resp)
		} else {
			return nil
		}
	case <-time.After(10 * time.Second):
		return fmt.Errorf("failed to set gateway deploy config: timeout")
	}
}
func (s *GatewayService) RemoveHost(ctx context.Context, hostid string) error {
	panic("implement me")
}
