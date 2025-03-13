package services

import (
	"context"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/model/request"
	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/proto"
)

type IGatewayService interface {
	AddHost(ctx context.Context, info *model.GatewayHost) error
	EditHost(ctx context.Context, hostid string, info *model.GatewayHost) error
	GetHost(ctx context.Context, hostid string) (model.GatewayHost, error)
	EditGateway(ctx context.Context, gatewayid, description string, tls *crypto.Tls) error
	DescriptionHost(ctx context.Context, hostid string) (*response.DescriptionHostResponse, error)
	DescriptionGateway(ctx context.Context, gatewayid string) (*response.DescriptionGatewayResponse, error)
	UserGetHosts(ctx context.Context, userid string) ([]model.GatewayHost, error)
	DeployGatewayInstance(ctx context.Context, hostid string, description string, tls *crypto.Tls) (string, error)

	GetRegisterStatus(_ context.Context, gatewayid string) *status.Status
	GetStatus(_ context.Context, gatewayid string) *status.Status

	GetErrorLogs(ctx context.Context, gatewayid string, from, to time.Time, limit int) ([]model.ErrorLogs, error)
	DescriptionError(ctx context.Context, errorID string) (model.ErrorLogs, error)

	HostGetGatewayDeployConfig(ctx context.Context, hostid string) (*model.ServerInfo, error)
	GatewayGetGatewayDeployConfig(ctx context.Context, gatewayid string) (*model.ServerInfo, error)

	HostSetGatewayDeployConfig(ctx context.Context, hostid string, info *model.ServerInfo) error
	GatewaySetGatewayDeployConfig(ctx context.Context, gatewayid string, info *model.ServerInfo) error

	// internal
	AddHostInternal(ctx context.Context, info *model.GatewayHost) error
	AddGatewayInternal(ctx context.Context, userID, gatewayID, description string, tls *crypto.Tls) error
	// AddAgentInternal add an agent to gateway(for internal called or debug), return agentID, and error
	AddAgentInternal(ctx context.Context, taskid *string, userid, gatewayid string, req *request.AddAgentRequest, openapidoc, enableFile []byte) (string, error)
	RemoveGatewayInternal(ctx context.Context, gatewayid string) error
	//RemoveGatewayHostInternal(ctx context.Context, hostid string) error

	// Task (for internal called or debug)

	// PushTask send task(the marshalled result) to gateway，
	//  return task-topic, taskID and error
	// if taskid is "", system will create a random taskid
	PushTask(ctx context.Context, taskID *string, gatewayID, userID string, bin []byte) (string, string, error)
}

type IAgentService interface {
	// PushTaskAgent send task(the marshalled result) to agent,
	//  return task-topic, taskID and error
	// if taskid is "", system will create a random taskid
	PushTaskAgent(_ context.Context, taskid *string, userID, gatewayID, agentID string, bin []byte) (string, string, error)
	// PushTaskAgentPb send task(the marshalled result) to agent, it like PushTaskAgent, but it will marshal pb to bin
	//  and pb type is core.TaskDetail
	PushTaskAgentPb(ctx context.Context, taskid *string, userID, gatewayID, agentID string, pb proto.Message) (string, string, error)

	GetAgentStatus(ctx context.Context, gatewayid string, ids []string) ([]model.AgentStatus, error)
	StartAgent(ctx context.Context, userid, gatewayid, agentid string) (taskid string, err error)
	StopAgent(ctx context.Context, userid, gatewayid, agentid string) (taskid string, err error)
	StartGather(ctx context.Context, userid, gatewayid, agentid string) (taskid string, err error)
	StopGather(ctx context.Context, userid, gatewayid, agentid string) (taskid string, err error)
	GetOpenApiDoc(ctx context.Context, userid, gatewayid, agentid string, docType agentpb.OpenapiDocType) (result *response.GetOpenApiDocResponse, err error)
}
