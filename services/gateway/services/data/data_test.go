package data

import (
	"math/rand"
	"testing"
	"time"

	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var timer = time.Now().Format("2006-01-02 15:04:05")

func newRandomDataMessage(n int) *agentpb.DataMessage {
	data := &agentpb.DataMessage{DataGatherStartTime: timer}

	var s [][]byte
	for i := 0; i < n; i++ {
		s = append(s, []byte(uuid.NewString()))
	}

	data.Data = s
	return data
}
func newTestAgent(n int) *data {
	agent := &data{data: make([]*agentpb.DataMessage, 0)}

	for i := 0; i < n; i++ {
		data := newRandomDataMessage(n)
		agent.parseData(data)
	}
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

	assert.Equal(randomTest*randomTest, int(data.DataLen))
	assert.Equal(timer, data.Time)
}
