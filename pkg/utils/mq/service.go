package mq

import (
	"context"
	"strings"
	"time"

	nats2 "github.com/AEnjoy/IoT-lubricant/pkg/utils/nats"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type Mq interface {
	Publish(topic string, msg any) error
	PublishBytes(topic string, msg []byte) error
	Subscribe(topic string) (<-chan any, error)
	Unsubscribe(topic string, sub <-chan any) error
	Close()
	GetPayLoad(sub <-chan any) any
	SetConditions(capacity int)
}

func NewMq() Mq {
	return &MessageQueue[any]{
		closeCh: make(chan struct{}),
	}
}

// NewNatsMq creates a new instance of NatsMq
func NewNatsMq[T any](url string) (*NatsMq, error) {
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

// NewKafkaMq creates a new instance of KafkaMq
func NewKafkaMq[T any](address, groupID string, partition, timeout int) *KafkaMq[T] {
	if timeout <= 0 {
		timeout = 10
	}
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
		timeout:     time.Duration(timeout) * time.Second,
		capacity:    100,
	}
}

// NewRedisMQ creates a new instance of Redis MQ
func NewRedisMQ[T any](addr string, password string, db int) (*RedisMq[T], error) {
	ctx, cancel := context.WithCancel(context.Background())
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		cancel()
		client.Close()
		return nil, err
	}

	return &RedisMq[T]{
		client:   client,
		ctx:      ctx,
		cancel:   cancel,
		channels: make(map[string]chan T),
		subs:     make(map[string]*redis.PubSub),
		capacity: 100,
	}, nil
}

// NewGoMq creates a new instance of Go internal datastruct implementation
func NewGoMq[T any]() *GoMq[T] {
	return &GoMq[T]{
		topics:   make(map[string][]chan T),
		capacity: 100,
	}
}
