package mq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMq_SetConditions(t *testing.T) {
	mq := NewMq()
	mq.SetConditions(5)

	sub, err := mq.Subscribe("testTopic")
	assert.NoError(t, err, "Subscribe should not return an error")

	err = mq.Publish("testTopic", "test message")
	assert.NoError(t, err, "Publish should not return an error")

	// 等待并检查消息是否能成功发布和接收
	select {
	case msg := <-sub:
		assert.Equal(t, "test message", msg, "Message received should match the message published")
	case <-time.After(3 * time.Second):
		t.Error("Message was not received in time")
	}

	// 清理
	assert.NoError(t, mq.Unsubscribe("testTopic", sub))
	mq.Close()
}

func TestMq_PublishAndSubscribe(t *testing.T) {
	mq := NewMq()
	mq.SetConditions(10)

	sub1, err1 := mq.Subscribe("topic1")
	sub2, err2 := mq.Subscribe("topic1")

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	err := mq.Publish("topic1", "Hello Subscribers")
	assert.NoError(t, err)

	select {
	case msg1 := <-sub1:
		assert.Equal(t, "Hello Subscribers", msg1)
	case <-time.After(3 * time.Second): // 超时时间小于3可能会触发超时，但实际场景中不会发生
		t.Error("Subscriber 1 did not receive the message in time")
	}

	select {
	case msg2 := <-sub2:
		assert.Equal(t, "Hello Subscribers", msg2)
	case <-time.After(1 * time.Second):
		t.Error("Subscriber 2 did not receive the message in time")
	}

	// Clear
	assert.NoError(t, mq.Unsubscribe("topic1", sub1))
	assert.NoError(t, mq.Unsubscribe("topic1", sub2))

}

func TestMq_Unsubscribe(t *testing.T) {
	mq := NewGoMq[string]()
	mq.SetConditions(10)

	sub, err := mq.Subscribe("topic1")
	assert.NoError(t, err)

	// 发布消息
	err = mq.Publish("topic1", "Message Before Unsubscribe")
	assert.NoError(t, err)

	// 确保订阅者收到消息
	select {
	case msg := <-sub:
		assert.Equal(t, "Message Before Unsubscribe", msg)
	case <-time.After(3 * time.Second):
		t.Error("Message not received before unsubscribe")
	}

	// 取消订阅
	err = mq.Unsubscribe("topic1", sub)
	assert.NoError(t, err)
	for range sub {
		// <-sub
		// drain the channel
	}

	// 再次发布消息，订阅者不应该收到任何消息
	err = mq.Publish("topic1", "Message After Unsubscribe")
	assert.NoError(t, err)

	select {
	case v := <-sub:
		if v != "" {
			t.Error("Unsubscribed subscriber should not receive any messages")
		}
	default:
		// Pass if no message is received
	}

	// 清理
	mq.Close()
}

func TestMq_Close(t *testing.T) {
	mq := NewMq()
	mq.SetConditions(10)

	sub, err := mq.Subscribe("topic1")
	assert.NoError(t, err)

	// 关闭消息队列
	mq.Close()

	// 确保已经关闭的队列不会再接收消息
	err = mq.Publish("topic1", "Message after close")
	assert.NoError(t, err)

	select {
	case _, ok := <-sub:
		assert.False(t, ok, "Channel should be closed after queue is closed")
	default:
		t.Error("Expected channel to be closed, but it is still open")
	}
}

func TestMq_GetPayLoad(t *testing.T) {
	mq := NewMq()
	mq.SetConditions(10)

	sub, err := mq.Subscribe("topic1")
	assert.NoError(t, err)

	// 发布消息
	err = mq.Publish("topic1", "Test Message")
	assert.NoError(t, err)

	// 使用 GetPayLoad 来获取消息
	payload := mq.GetPayLoad(sub)
	assert.Equal(t, "Test Message", payload)

	// 清理
	assert.NoError(t, mq.Unsubscribe("topic1", sub))
	mq.Close()
}
