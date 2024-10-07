package mq

import (
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/assert"
)

func TestNatsMq(t *testing.T) {
	assert := assert.New(t)
	opts := &server.Options{
		Port: 4223,
	}

	natsServer, err := server.NewServer(opts)
	assert.NoError(err)

	t.Log("starting nats server")
	go natsServer.Start()
	if !natsServer.ReadyForConnections(10 * time.Second) {
		t.Fatal("nats server did not start")
	}
	defer natsServer.Shutdown()

	t.Log("starting mq")
	mq, err := NewNatsMq[string]("nats://127.0.0.1:4223")
	assert.NoError(err)
	defer mq.Close()
	mq.SetConditions(100)

	t.Log("Test: Publish message(with no subscribed)")
	assert.NoError(mq.Publish("test", "hello"))

	t.Log("Test: Subscribe message")
	ch, err := mq.Subscribe("test")
	assert.NoError(err)

	t.Log("Test: Publish message(with subscribed)")
	assert.NoError(mq.Publish("test", "hello"))

	t.Log("get payload")
	payload := mq.GetPayLoad(ch)
	assert.Equal("hello", payload)

	t.Log("unsubscribe message")
	assert.NoError(mq.Unsubscribe("test", ch))
}
