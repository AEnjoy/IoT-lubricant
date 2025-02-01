package mq

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKafkaMq(t *testing.T) {
	address := os.Getenv("KAFKA_BROKE")
	if len(address) == 0 {
		t.Skip("This unit test will only be executed when and only when the environment variable KAFKA_BROKE is set")
	}
	assert := assert.New(t)

	mq := NewKafkaMq[string](address, "", 0, 10)
	defer mq.Close()

	mq.SetConditions(10)

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

func BenchmarkKafkaMq(b *testing.B) {
	address := os.Getenv("KAFKA_BROKE")
	if len(address) == 0 {
		b.Skip("This unit test will only be executed when and only when the environment variable KAFKA_BROKE is set")
	}
	mq := NewKafkaMq[string](address, "", 0, 10)
	defer mq.Close()

	mq.SetConditions(10)
	testBenchmark(b, mq)
}
