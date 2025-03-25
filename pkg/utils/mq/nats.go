//go:build !jetstream

package mq

import (
	"sync"

	nats2 "github.com/aenjoy/iot-lubricant/pkg/utils/nats"
	json "github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
)

type NatsMq struct {
	nc       *nats.Conn
	subs     map[string]*nats.Subscription
	channels map[string]chan any
	mu       sync.Mutex
	capacity int
}

// Publish sends a message to the specified topic
func (mq *NatsMq) Publish(topic string, msg any) error {
	return mq.nc.Publish(topic, msgToBytes(msg)) // Helper function to convert message to []byte
}
func (mq *NatsMq) PublishBytes(topic string, msg []byte) error {
	return mq.nc.Publish(topic, msg) // Helper function to convert message to []byte
}

// Subscribe subscribes to a topic and returns a channel for receiving messages
func (mq *NatsMq) Subscribe(topic string) (<-chan any, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if _, exists := mq.channels[topic]; exists {
		return nil, ErrHasBeenSubscribed // Already subscribed
	}

	ch := make(chan any, mq.capacity) // Create a buffered channel
	sub, err := mq.nc.Subscribe(topic, func(msg *nats.Msg) {
		var data any
		bytesToMsg(msg.Data, &data) // Helper function to convert []byte to T
		ch <- data
	})
	if err != nil {
		return nil, err
	}

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
	sub, err := mq.nc.Subscribe(topic, func(msg *nats.Msg) {
		ch <- msg.Data
	})
	if err != nil {
		return nil, err
	}

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
func NewNatsMq(url string) (*NatsMq, error) {
	nc, err := nats2.NewNatsClient(url)
	if err != nil {
		return nil, err
	}
	return &NatsMq{
		nc:       nc,
		subs:     make(map[string]*nats.Subscription),
		channels: make(map[string]chan any),
		capacity: 10, // default capacity
	}, nil
}
