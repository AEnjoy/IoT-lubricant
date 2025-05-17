package rpcserver

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/types"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
	"github.com/aenjoy/iot-lubricant/services/corepkg/config"
	"github.com/aenjoy/iot-lubricant/services/corepkg/dataapi"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

var ErrProjectNotFound = errors.New("project not found")

func (*PbCoreServiceImpl) taskSendErrorMessage(s grpc.BidiStreamingServer[corepb.Task, corepb.Task], code int, message string) {
	var out corepb.Task
	var errorResp corepb.Task_ErrorMessage
	errorResp.ErrorMessage = &metapb.ErrorMessage{Code: &status.Status{Code: int32(code), Message: message}}
	out.Task = &errorResp
	_ = s.Send(&out)
}
func (*PbCoreServiceImpl) gatewayOffline(mq mq.Mq, userid, gatewayid string) error {
	logg.L.Debugf("gateway offline: %s", gatewayid)
	return mq.PublishBytes(fmt.Sprintf("/monitor/%s/offline", taskTypes.TargetGateway),
		[]byte(fmt.Sprintf("%s<!SPLIT!>%s", userid, gatewayid)))
}

// getGatewayID 获取调用接口的网关ID
//
//	todo:(现阶段其实应该叫做网关名,后续需要使用真正的id)
func (*PbCoreServiceImpl) getGatewayID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata not found")
	}

	gatewayIDs := md.Get(string(types.NameGatewayID))
	if len(gatewayIDs) == 0 {
		return "", fmt.Errorf("gatewayid not found in metadata")
	}

	return gatewayIDs[0], nil
}
func (*PbCoreServiceImpl) getUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata not found")
	}

	userIDs := md.Get(constant.USER_ID)
	if len(userIDs) == 0 {
		return "", fmt.Errorf("userid not found in metadata")
	}

	return userIDs[0], nil
}

var (
	_initDataStoreServiceClientMutex sync.Mutex
	_storeDataCliStreamMutex         sync.Mutex

	serv               *dataapi.DataStoreApiService
	storeDataCliStream grpc.ClientStreamingClient[svcpb.StoreDataRequest, svcpb.StoreDataResponse]
)

func (i *PbCoreServiceImpl) _svcDataStoreServiceClientDataStream(ctx context.Context, data []byte, projectId string) error {
	_storeDataCliStreamMutex.Lock()
	if storeDataCliStream == nil || storeDataCliStream.Context().Err() != nil {
		storeDataCli, err := serv.StoreData(ctx)
		if err != nil {
			logg.L.Errorf("failed to create store data stream: %v", err)
			_storeDataCliStreamMutex.Unlock()
			return err
		}
		storeDataCliStream = storeDataCli
	}
	_storeDataCliStreamMutex.Unlock()

	err := storeDataCliStream.Send(&svcpb.StoreDataRequest{
		Data:      data,
		ProjectID: projectId,
	})
	if err != nil {
		logg.L.Errorf("failed to send data to store data stream: %v", err)
	}
	return err
}
func (i *PbCoreServiceImpl) handelRecvData(ctx context.Context, data *corepb.Data, projectId string) {
	err := i.pool.Submit(func() {
		dataCli := i.DataStore
		marshal, err := proto.Marshal(data)
		if err != nil {
			logg.L.Errorf("failed to marshal data: %v", err)
		}
		if config.GetConfig().SvcDataStoreMode == "rpc" {
			_initDataStoreServiceClientMutex.Lock()
			if serv == nil {
				serv = ioc.Controller.Get(ioc.APP_NAME_CORE_Internal_DATASTORE_API_SERVICE).(*dataapi.DataStoreApiService)
			}
			_initDataStoreServiceClientMutex.Unlock()
			err := i._svcDataStoreServiceClientDataStream(ctx, marshal, projectId)
			if err != nil {
				logg.L.Errorf("failed to send data to datastore service: %v", err)
				return
			}
			return
		}
		err = dataCli.V2mq.Publish(constant.DATASTORE_PROJECT, []byte(projectId))
		if err != nil {
			logg.L.Errorf("failed to publish projectId[%s] to datastore topic: %v", projectId, err)
			return
		}

		err = dataCli.V2mq.QueuePublish(fmt.Sprintf(constant.DATASTORE_PROJECT_DATA, projectId), marshal)
		if err != nil {
			logg.L.Errorf("failed to publish data: %v", err)
		}
	})
	if err != nil {
		logg.L.Errorf("failed to create save data to store task thread")
	}
}
func (i *PbCoreServiceImpl) getProjectId(ctx context.Context, agentid string) (string, error) {
	id, err := i.CacheCli.Get(ctx, fmt.Sprintf("project-id-agent-%s", agentid))
	if err != nil || id == "" {
		i.getProjectIdMutex.Lock()
		defer i.getProjectIdMutex.Unlock()

		project, _ := i.DataStore.GetProjectByAgentID(ctx, agentid)
		if project.ProjectID == "" {
			return "", ErrProjectNotFound
		}
		id = project.ProjectID
		err = i.CacheCli.Set(ctx, fmt.Sprintf("project-id-agent-%s", agentid), id)
	}
	return id, err
}
