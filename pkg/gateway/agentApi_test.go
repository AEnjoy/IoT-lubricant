package gateway

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	grpcmock "github.com/AEnjoy/IoT-lubricant/pkg/mock/grpc"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAPP_JoinAgent(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	mockMqClient := mq.NewMockMq[[]byte](ctrl)
	mockGrpcClient := grpcmock.NewCoreServiceClient(t)
	mockGrpcDataStream := grpcmock.NewBidiStreamingServer[core.Data, core.Data](t)
	deviceList := &sync.Map{}

	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(8*time.Second))
	defer cf()

	var success bool

	app := &app{
		mq:         mockMqClient,
		ctrl:       ctx,
		deviceList: deviceList,
		grpcClient: mockGrpcClient,
		clientMq: &clientMq{
			ctrl:       ctx,
			cancel:     cf,
			deviceList: deviceList,
		},
	}
	mockGrpcClient.On("PushData", mock.Anything).Return(mockGrpcDataStream, nil)
	for {
		select {
		case <-ctx.Done():
			assert.True(success)
			return
		case <-time.Tick(time.Second):
			id := uuid.NewString()
			regCh := make(chan []byte)
			nilCh := make(chan []byte) // agentDevice dataPush messagePush

			ping, err := json.Marshal(model.Ping{Status: 1})
			assert.NoError(err)

			// Subscribe Topic_AgentRegister
			mockMqClient.EXPECT().Subscribe(model.Topic_AgentRegister+id).Return(regCh, nil)

			// handelAgentRegister
			// send register success to channel
			go func() {
				t.Logf("Test: Join %s to Gateway", id)
				data, err := json.Marshal(model.Register{ID: id})
				assert.NoError(err)
				regCh <- data
			}()
			mockMqClient.EXPECT().Publish(model.Topic_AgentRegisterAck+id, ping).Return(nil)

			// subscribeDeviceMQ
			mockMqClient.EXPECT().Subscribe(model.Topic_AgentDevice+id).Return(nilCh, nil)

			// handelAgentDataPush
			// Subscribe Topic_AgentDataPush
			mockMqClient.EXPECT().Subscribe(model.Topic_AgentDataPush+id).Return(nilCh, nil)

			// Subscribe Topic_MessagePush
			mockMqClient.EXPECT().Subscribe(model.Topic_MessagePush+id).Return(nilCh, nil)

			assert.NoError(app.joinAgent(id))
			success = true
		}
	}
}
