package mq

import (
	nats2 "github.com/AEnjoy/IoT-lubricant/pkg/utils/nats"
	"github.com/nats-io/nats.go"
)

type Mq[T any] interface {
	Publish(topic string, msg T) error
	Subscribe(topic string) (<-chan T, error)
	Unsubscribe(topic string, sub <-chan T) error
	Close()
	GetPayLoad(sub <-chan T) T
	SetConditions(capacity int)
}

func NewMq[T any]() Mq[T] {
	return &MessageQueue[T]{
		closeCh: make(chan struct{}),
	}
}

// NewNatsMq creates a new instance of NatsMq
func NewNatsMq[T any](url string) (*NatsMq[T], error) {
	nc, err := nats2.NewNatsClient(url)
	if err != nil {
		return nil, err
	}
	return &NatsMq[T]{
		nc:       nc,
		subs:     make(map[string]*nats.Subscription),
		channels: make(map[string]chan T),
		capacity: 10, // default capacity
	}, nil
}
