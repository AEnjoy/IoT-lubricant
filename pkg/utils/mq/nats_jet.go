//go:build jetstream

/*
todo: 经过测试,目前 jetstream 存在bug:发布消息timeout,需要修复
*/
package mq

import (
	"sync"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
)

type NatsMq struct {
	nc       *nats.Conn
	js       nats.JetStreamContext
	subs     map[string]*nats.Subscription
	channels map[string]chan any
	mu       sync.Mutex
	capacity int
}

// Publish sends a message to the specified topic
func (mq *NatsMq) Publish(topic string, msg any) error {
	_, err := mq.js.Publish(topic, msgToBytes(msg)) // Helper function to convert message to []byte
	return err
}
func (mq *NatsMq) PublishBytes(topic string, msg []byte) error {
	_, err := mq.js.Publish(topic, msg) // Helper function to convert message to []byte
	return err
}

// Subscribe subscribes to a topic and returns a channel for receiving messages
func (mq *NatsMq) Subscribe(topic string) (<-chan any, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if _, exists := mq.channels[topic]; exists {
		return nil, ErrHasBeenSubscribed // Already subscribed
	}

	ch := make(chan any, mq.capacity) // Create a buffered channel
	sub, err := mq.js.PullSubscribe(topic, "IUTER_CONSUMER",
		nats.DeliverAll(),
		nats.AckExplicit(),
	)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			msgs, _ := sub.Fetch(10, nats.MaxWait(3*time.Second))
			for _, msg := range msgs {
				bytesToMsg(msg.Data, &ch)
				_ = msg.Ack()
			}
		}
	}()

	mq.subs[topic] = sub
	mq.channels[topic] = ch
	return ch, nil
}
func (mq *NatsMq) SubscribeBytes(topic string) (<-chan any, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if _, exists := mq.channels[topic]; exists {
		return nil, ErrHasBeenSubscribed // Already subscribed
	}

	ch := make(chan any, mq.capacity) // Create a buffered channel
	sub, err := mq.js.PullSubscribe(topic, "IUTER_CONSUMER",
		nats.DeliverAll(),
		nats.AckExplicit(),
	)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			msgs, _ := sub.Fetch(10, nats.MaxWait(3*time.Second))
			for _, msg := range msgs {
				ch <- msg.Data
				_ = msg.Ack()
			}
		}
	}()

	mq.subs[topic] = sub
	mq.channels[topic] = ch
	return ch, nil
}

// Unsubscribe cancels the subscription and closes the channel
func (mq *NatsMq) Unsubscribe(topic string, sub <-chan any) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if subscription, exists := mq.subs[topic]; exists {
		_ = subscription.Unsubscribe() // Unsubscribe from NATS
		close(mq.channels[topic])      // Close the channel
		delete(mq.subs, topic)         // Remove from map
		delete(mq.channels, topic)
		return nil
	}
	return ErrNotFoundSubscriber // Subscription not found
}

// Close closes the NATS connection and cleans up resources
func (mq *NatsMq) Close() {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for topic, subscription := range mq.subs {
		_ = subscription.Unsubscribe()
		close(mq.channels[topic])
	}
	mq.nc.Close()
}

// GetPayLoad gets the payload from the subscription channel
func (mq *NatsMq) GetPayLoad(sub <-chan any) any {
	return <-sub
}

// SetConditions sets the capacity of the channel
func (mq *NatsMq) SetConditions(capacity int) {
	mq.capacity = capacity
}

// Helper function to convert message to bytes (T -> []byte)
func msgToBytes[T any](msg T) []byte {
	// Implement serialization logic (e.g., JSON encoding)
	// This will depend on the type of T
	// Example using JSON:
	data, _ := json.Marshal(msg)
	return data
}

// Helper function to convert bytes to message ([]byte -> T)
func bytesToMsg[T any](data []byte, msg *T) {
	// Implement deserialization logic (e.g., JSON decoding)
	// Example using JSON:
	_ = json.Unmarshal(data, msg)
}

// NewNatsMq creates a new instance of NatsMq
func NewNatsMq[T any](url string) (*NatsMq, error) {
	nc, err := nats2.NewNatsClient(url)
	if err != nil {
		return nil, err
	}
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "IUTER_STREAM",
		Subjects: []string{">"},
		Storage:  nats.MemoryStorage,
		MaxAge:   1 * time.Hour,
		NoAck:    true,
	})
	if err != nil {
		return nil, err
	}
	return &NatsMq{
		nc:       nc,
		js:       js,
		subs:     make(map[string]*nats.Subscription),
		channels: make(map[string]chan any),
		capacity: 10, // default capacity
	}, nil
}
