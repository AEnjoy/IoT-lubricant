package gateway

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/compress"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var timer = time.Now().Format("2006-01-02 15:04:05")

func newRandomDataMessage(n int) *gateway.DataMessage {
	data := &gateway.DataMessage{
		Flag:      2,
		MessageId: uuid.NewString(),
		AgentId:   uuid.NewString(),
		Time:      timer,
	}

	var s [][]byte
	for i := 0; i < n; i++ {
		s = append(s, []byte(uuid.NewString()))
	}

	data.Data = bytes.Join(s, compress.Sepa)
	return data
}
func newTestAgent(n int) *agentData {
	agent := &agentData{data: make([]*gateway.DataMessage, 0)}

	data := newRandomDataMessage(n)
	agent.parseData(data, 1)
	return agent
}
func TestParseData(t *testing.T) {
	assert := assert.New(t)
	randomTest := rand.Intn(71) + 30 // 30-100
	agent := newTestAgent(randomTest)

	assert.Equal(randomTest, len(agent.data))
}
func TestCoverToGrpcData(t *testing.T) {
	assert := assert.New(t)
	randomTest := rand.Intn(71) + 30 // 30-100
	data := newTestAgent(randomTest).coverToGrpcData()

	assert.Equal(randomTest, int(data.DataLen))
	assert.Equal(timer, data.Time)
}
