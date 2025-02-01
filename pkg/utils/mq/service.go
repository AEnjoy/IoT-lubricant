package mq

import (
	"strings"
	"time"

	nats2 "github.com/AEnjoy/IoT-lubricant/pkg/utils/nats"
	"github.com/nats-io/nats.go"
	"github.com/segmentio/kafka-go"
)

type Mq[T any] interface {
	Publish(topic string, msg T) error
	PublishBytes(topic string, msg []byte) error
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

// NewKafkaMq creates a new instance of KafkaMq
func NewKafkaMq[T any](address, groupID string, partition, timeout int) *KafkaMq[T] {
	brokers := strings.Split(address, ",")
	adminClient := &kafka.Client{
		Addr:    kafka.TCP(brokers...),
		Timeout: time.Duration(timeout) * time.Second,
	}

	return &KafkaMq[T]{
		address:     brokers,
		groupID:     groupID,
		partition:   partition,
		adminClient: adminClient,

		writers:     make(map[string]*kafka.Writer),
		subscribers: make(map[string][]*subscriber[T]),
		timeout:     10 * time.Second,
		capacity:    100,
	}
}
