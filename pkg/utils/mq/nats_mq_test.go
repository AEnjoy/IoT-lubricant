package mq

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNatsMq(t *testing.T) {
	assert := assert.New(t)
	natsServer, url := natsUrl(t)
	defer natsServer()

	mq, err := NewNatsMq(url)
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

var subTestResult sync.Map

func consumerNatsClient(t *testing.T, url string, i int, ctx context.Context) {
	assert := assert.New(t)
	cli, err := NewNatsMq(url)
	assert.NoError(err)

	subscribe, err := cli.Subscribe("/test/topic/123")
	assert.NoError(err)

	for {
		select {
		case <-ctx.Done():
			return
		case <-subscribe:
			go subTestResult.Store(i, struct{}{})
		}
	}
}
