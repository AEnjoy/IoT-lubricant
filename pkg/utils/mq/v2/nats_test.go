//go:build !jetstream

package mq

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func natsUrl(t *testing.T) (func(), string) {
	url, ok := os.LookupEnv("TEST_NATS_MQ_URL")
	if !ok {
		t.Log("env TEST_NATS_MQ_URL not found, use internal nats")
		opts := &server.Options{
			Port:    4222 + rand.Intn(100),
			Debug:   true,
			Trace:   true,
			LogFile: "nats.log",
		}
		natsServer, err := server.NewServer(opts)
		assert.NoError(t, err)

		t.Log("starting nats server")
		go natsServer.Start()
		if !natsServer.ReadyForConnections(10 * time.Second) {
			t.Fatal("nats server did not start")
		}

		url = natsServer.ClientURL()
		return func() {
			t.Log("shutting down nats server")
			natsServer.Shutdown()
		}, url
	}
	return func() {}, url
}

func TestNatsMqAbility(t *testing.T) {
	assert := assert.New(t)
	clean, url := natsUrl(t)
	defer clean()

	mq, err := NewNatsMq(url)
	assert.NoError(err)
	defer mq.Close()
	mq.SetConditions(100)
	t.Run("Production and re consumption", func(t *testing.T) {
		t.Skip("Only in JetStream Enabled can pass this test")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for i := 0; i < 100; i++ {
			assert.NoError(mq.Publish("/test/topic/123", []byte("hello")))
		}
		time.Sleep(time.Second)
		for i := 0; i < 10; i++ {
			go consumerNatsClient(t, url, i, ctx)
		}
		time.Sleep(3 * time.Second)

		var recv = 0
		subTestResult.Range(func(key, value any) bool {
			assert.Equal(value, struct{}{})
			recv++
			return true
		})
		assert.Equal(10, recv)
	})

	subTestResult = sync.Map{}
	t.Run("Consumption regeneration", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for i := 0; i < 10; i++ {
			go consumerNatsClient(t, url, i, ctx)
		}
		time.Sleep(time.Second)
		for i := 0; i < 100; i++ {
			assert.NoError(mq.Publish("/test/topic/123", []byte("hello")))
		}
		time.Sleep(3 * time.Second)
		var recv = 0
		subTestResult.Range(func(key, value any) bool {
			assert.Equal(value, struct{}{})
			recv++
			return true
		})
		assert.Equal(10, recv)
	})
}

func TestNatsSubscribe(t *testing.T) {
	clean, url := natsUrl(t)
	defer clean()

	mq, err := NewNatsMq(url)
	assert.NoError(t, err)
	defer mq.Close()
	mq.SetConditions(10) // 设置较小缓冲便于测试溢出情况

	t.Run("ReceiveMq", func(t *testing.T) {
		topic := "test.sub.normal"
		recvCh, err := mq.Subscribe(topic)
		assert.NoError(t, err)

		msg := []byte("normal_msg")
		assert.NoError(t, mq.Publish(topic, msg))

		select {
		case received := <-recvCh:
			assert.Equal(t, msg, received)
		case <-time.After(1 * time.Second):
			t.Fatal("TimeOut")
		}
	})

	t.Run("Not allow duplicate subscribe", func(t *testing.T) {
		topic := "test.duplicate"
		_, err := mq.Subscribe(topic)
		assert.NoError(t, err)

		_, err = mq.Subscribe(topic)
		assert.ErrorContains(t, err, "already subscribed")
	})

	t.Run("Drop Message", func(t *testing.T) {
		topic := "test.buffer.full"
		mq.SetConditions(2)
		recvCh, err := mq.Subscribe(topic)
		assert.NoError(t, err)

		for i := 0; i < 5; i++ {
			mq.Publish(topic, []byte(fmt.Sprintf("msg%d", i)))
		}

		time.Sleep(500 * time.Millisecond)
		assert.Len(t, recvCh, 2)
	})
}
func TestNatsMq_QueueSubscribe_LoadBalancing(t *testing.T) {
	// Setup: Start NATS Server using the helper
	shutdownNats, natsURL := natsUrl(t)
	defer shutdownNats() // Ensure NATS server is stopped after test

	// Setup: Create two separate MQ instances (simulating different consumers/processes)
	mq1, err := NewNatsMq(natsURL)
	require.NoError(t, err, "Failed to create NatsMq instance 1")
	defer mq1.Close()

	mq2, err := NewNatsMq(natsURL)
	require.NoError(t, err, "Failed to create NatsMq instance 2")
	defer mq2.Close()

	// --- Test Parameters ---
	queueTopic := "work.tasks.test"
	numMessages := 10 // Number of messages to publish

	// --- Subscribe from both instances to the same queue topic ---
	// They should join the same queue group derived from queueTopic
	qChan1, err := mq1.QueueSubscribe(queueTopic)
	require.NoError(t, err, "mq1 failed to queue subscribe")
	t.Logf("MQ1 subscribed to queue topic: %s", queueTopic)

	qChan2, err := mq2.QueueSubscribe(queueTopic)
	require.NoError(t, err, "mq2 failed to queue subscribe")
	t.Logf("MQ2 subscribed to queue topic: %s", queueTopic)

	// --- Goroutines to collect messages from each subscriber ---
	var mu sync.Mutex // Mutex to protect access to received maps

	// Map to store messages received by mq1: map[message_content]count
	received1 := make(map[string]int)
	// Map to store messages received by mq2: map[message_content]count
	received2 := make(map[string]int)

	processMessages := func(id int, ch <-chan any, receivedMap map[string]int) {
		t.Logf("Worker %d started, listening...", id)
		for msgAny := range ch {
			// Assert that received message is []byte and convert to string
			msgData, ok := msgAny.([]byte)
			require.True(t, ok, "Worker %d received non-byte data: %T", id, msgAny)
			msgStr := string(msgData)

			mu.Lock()
			receivedMap[msgStr]++ // Increment count for this message
			t.Logf("Worker %d received: %s (Count: %d)", id, msgStr, receivedMap[msgStr])
			mu.Unlock()
		}
	}

	go processMessages(1, qChan1, received1)
	go processMessages(2, qChan2, received2)

	// --- Publish messages ---
	// Allow subscribers a moment to be ready
	time.Sleep(200 * time.Millisecond)
	t.Logf("Publishing %d messages to topic: %s", numMessages, queueTopic)
	publishedMessages := make(map[string]bool) // Keep track of what was published
	for i := 0; i < numMessages; i++ {
		message := fmt.Sprintf("message-%d", i)
		publishedMessages[message] = true
		err := mq1.QueuePublish(queueTopic, []byte(message)) // Use either mq1 or mq2 to publish
		require.NoError(t, err, "Failed to publish message %d", i)
	}
	t.Log("Finished publishing messages.")

	// --- Wait briefly for messages to be processed ---
	// In a real system, you'd wait for confirmation or use other sync mechanisms.
	// Here, we give NATS some time to distribute.
	time.Sleep(500 * time.Millisecond)

	// --- Unsubscribe and wait for goroutines to finish ---
	t.Log("Unsubscribing worker 1...")
	err = mq1.QueueUnsubscribe(queueTopic)
	require.NoError(t, err, "Failed to unsubscribe mq1")

	t.Log("Unsubscribing worker 2...")
	err = mq2.QueueUnsubscribe(queueTopic)
	require.NoError(t, err, "Failed to unsubscribe mq2")

	t.Log("Worker goroutines finished.")

	// --- Assertions ---
	mu.Lock() // Lock for final assertion checks
	defer mu.Unlock()

	totalReceivedCount := 0
	uniqueReceivedMessages := make(map[string]bool)

	t.Logf("Messages received by Worker 1: %d", len(received1))
	for msg, count := range received1 {
		assert.Equal(t, 1, count, "Worker 1 received message '%s' %d times, expected 1", msg, count)
		totalReceivedCount += count
		uniqueReceivedMessages[msg] = true
	}

	t.Logf("Messages received by Worker 2: %d", len(received2))
	for msg, count := range received2 {
		assert.Equal(t, 1, count, "Worker 2 received message '%s' %d times, expected 1", msg, count)
		totalReceivedCount += count
		uniqueReceivedMessages[msg] = true
		// CRITICAL CHECK: Assert that a message received by worker 2 was NOT received by worker 1
		_, foundIn1 := received1[msg]
		assert.False(t, foundIn1, "Message '%s' was received by BOTH workers!", msg)
	}

	// Check total messages received vs published
	assert.Equal(t, numMessages, totalReceivedCount, "Total messages received across workers (%d) does not match number published (%d)", totalReceivedCount, numMessages)

	// Check if all published messages were received by at least one worker
	assert.Equal(t, len(publishedMessages), len(uniqueReceivedMessages), "Number of unique received messages (%d) does not match number of unique published messages (%d)", len(uniqueReceivedMessages), len(publishedMessages))
	for pubMsg := range publishedMessages {
		_, found := uniqueReceivedMessages[pubMsg]
		assert.True(t, found, "Published message '%s' was not received by any worker", pubMsg)
	}

	t.Log("Assertions complete.")
}
