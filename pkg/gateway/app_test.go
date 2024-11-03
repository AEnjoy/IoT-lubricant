package gateway

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/mock/db"
	grpcmock "github.com/AEnjoy/IoT-lubricant/pkg/mock/grpc"
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const TestTime = 8

func TestGatewayAPP(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)

	mockMqClient := mq.NewMockMq[[]byte](ctrl)
	mockDbClient := db.NewMockGatewayDbOperator(ctrl)
	mockGrpcClient := grpcmock.NewCoreServiceClient(t)
	mockGrpcTaskStream := grpcmock.NewBidiStreamingServer[core.Task, core.Task](t)
	mockGrpcDataStream := grpcmock.NewBidiStreamingServer[core.Data, core.Data](t)

	deviceList := &sync.Map{}

	ctx, cf := context.WithDeadline(
		context.WithValue(context.Background(), types.NameGatewayID, uuid.NewString()),
		time.Now().Add(TestTime*time.Second))
	defer cf()

	app := &app{
		mq:                mockMqClient,
		ctrl:              ctx,
		deviceList:        deviceList,
		GatewayDbOperator: mockDbClient,
		grpcClient:        mockGrpcClient,
	}

	// mock link an agent to gateway
	agentId := uuid.NewString()
	nilCh := make(chan []byte) // agentDevice dataPush messagePush
	regCh := make(chan []byte, 1)
	dataPushCh := make(chan []byte, 1)

	mockDbClient.EXPECT().GetAllAgentId().Return([]string{agentId})
	mockDbClient.EXPECT().GetAgentGatherCycle(agentId).Return(1)
	mockDbClient.EXPECT().GetAllAgents().Return([]types.Agent{}, nil)

	mockMqClient.EXPECT().Subscribe(types.Topic_AgentRegister+agentId).Return(regCh, nil)

	// send register ack to agent
	ping, err := json.Marshal(types.Ping{Status: 1})
	assert.NoError(err)
	mockMqClient.EXPECT().Publish(types.Topic_AgentRegisterAck+agentId, ping).Return(nil)

	// mock-agent register
	data, err := json.Marshal(types.Register{ID: agentId})
	assert.NoError(err)
	regCh <- data

	// mock agent register and data send
	mockMqClient.EXPECT().Subscribe(types.Topic_AgentDataPush+agentId).Return(dataPushCh, nil)
	messageData, err := json.Marshal(newRandomDataMessage(45))
	assert.NoError(err)
	dataPushCh <- messageData

	// mock grpc server
	mockGrpcClient.On("GetTask", mock.Anything).Return(mockGrpcTaskStream, nil)
	mockGrpcClient.On("PushData", mock.Anything).Return(mockGrpcDataStream, nil)

	var (
		resp                   core.Task
		respDetail             core.TaskDetail
		respGetTask            core.Task_GatewayGetTaskResponse
		GatewayGetTaskResponse core.GatewayGetTaskResponse
		command                types.TaskCommand
	)

	data, err = json.Marshal(command)
	assert.NoError(err)
	respDetail.Content = data
	GatewayGetTaskResponse.Message = &respDetail
	respGetTask.GatewayGetTaskResponse = &GatewayGetTaskResponse
	resp.Task = &respGetTask

	mockGrpcTaskStream.On("Recv").Return(&resp, nil)
	mockGrpcDataStream.On("Recv").WaitUntil(time.Tick((TestTime+4)*time.Second)).Return(&resp, nil)
	mockGrpcDataStream.On("Send", mock.Anything).Return(nil)

	// mock mq Subscribe and publish
	{
		// handelGatewayInfo
		mockMqClient.EXPECT().Subscribe(types.Topic_GatewayInfo).Return(nilCh, nil)

		// handelGatewayData
		mockMqClient.EXPECT().Subscribe(types.Topic_GatewayData).Return(nilCh, nil)

		// handelPing
		mockMqClient.EXPECT().Subscribe(types.Topic_Ping).Return(nilCh, nil)

		// subscribeDeviceMQ
		mockMqClient.EXPECT().Subscribe(types.Topic_AgentDevice+agentId).Return(nilCh, nil)

		// Subscribe Topic_MessagePush
		mockMqClient.EXPECT().Subscribe(types.Topic_MessagePush+agentId).Return(nilCh, nil)
	}
	go func() {
		assert.NoError(app.Run())
	}()

	<-time.After(TestTime * time.Second)
}
