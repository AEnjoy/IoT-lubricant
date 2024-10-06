package mq

import (
	"encoding/json"
	"sync"

	"github.com/nats-io/nats.go"
)

type NatsMq[T any] struct {
	nc       *nats.Conn
	subs     map[string]*nats.Subscription
	channels map[string]chan T
	mu       sync.Mutex
	capacity int
}

// Publish sends a message to the specified topic
func (mq *NatsMq[T]) Publish(topic string, msg T) error {
	return mq.nc.Publish(topic, msgToBytes(msg)) // Helper function to convert message to []byte
}

// Subscribe subscribes to a topic and returns a channel for receiving messages
func (mq *NatsMq[T]) Subscribe(topic string) (<-chan T, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if _, exists := mq.channels[topic]; exists {
		return nil, ErrHasBeenSubscribed // Already subscribed
	}

	ch := make(chan T, mq.capacity) // Create a buffered channel
	sub, err := mq.nc.Subscribe(topic, func(msg *nats.Msg) {
		var data T
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

// Unsubscribe cancels the subscription and closes the channel
func (mq *NatsMq[T]) Unsubscribe(topic string, sub <-chan T) error {
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
func (mq *NatsMq[T]) Close() {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for topic, subscription := range mq.subs {
		_ = subscription.Unsubscribe()
		close(mq.channels[topic])
	}
	mq.nc.Close()
}

// GetPayLoad gets the payload from the subscription channel
func (mq *NatsMq[T]) GetPayLoad(sub <-chan T) T {
	return <-sub
}

// SetConditions sets the capacity of the channel
func (mq *NatsMq[T]) SetConditions(capacity int) {
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
