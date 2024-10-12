package gateway

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/mock/db"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// grpc Mock
// https://github.com/nhatthm/grpcmock
func RegisterServiceServer(s grpc.ServiceRegistrar, srv core.CoreServiceServer) {
	s.RegisterService(&core.CoreService_ServiceDesc, srv)
}
func TestGatewayAPP(t *testing.T) {
	t.Skip("not all implement yet")
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	mockMqClient := mq.NewMockMq[[]byte](ctrl)
	mockDbClient := db.NewMockGatewayDbCli(ctrl)
	mockGrpcClient := NewMockCoreServiceClient(ctrl)

	deviceList := &sync.Map{}

	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(8*time.Second))
	defer cf()

	app := &app{
		mq:           mockMqClient,
		ctrl:         ctx,
		deviceList:   deviceList,
		GatewayDbCli: mockDbClient,
		grpcClient:   mockGrpcClient,
		clientMq: &clientMq{
			ctrl:       ctx,
			cancel:     cf,
			deviceList: deviceList,
		},
	}
	// mock link an agent to gateway
	agentId := uuid.NewString()
	regCh := make(chan []byte)
	nilCh := make(chan []byte) // agentDevice dataPush messagePush

	ping, err := json.Marshal(model.Ping{Status: 1})
	assert.NoError(err)

	// Subscribe Topic_AgentRegister
	mockMqClient.EXPECT().Subscribe(model.Topic_AgentRegister+agentId).Return(regCh, nil)

	// handelAgentRegister
	// send register success to channel
	go func() {
		t.Logf("Join %s to Gateway", agentId)
		data, err := json.Marshal(model.Register{ID: agentId})
		assert.NoError(err)
		regCh <- data
	}()
	mockMqClient.EXPECT().Publish(model.Topic_AgentRegisterAck+agentId, ping).Return(nil)

	// subscribeDeviceMQ
	mockMqClient.EXPECT().Subscribe(model.Topic_AgentDevice+agentId).Return(nilCh, nil)

	// handelAgentDataPush
	// Subscribe Topic_AgentDataPush
	mockMqClient.EXPECT().Subscribe(model.Topic_AgentDataPush+agentId).Return(nilCh, nil)

	// Subscribe Topic_MessagePush
	mockMqClient.EXPECT().Subscribe(model.Topic_MessagePush+agentId).Return(nilCh, nil)
	go func() {
		assert.NoError(app.Run())
	}()

}
