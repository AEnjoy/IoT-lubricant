package mq

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type Mq interface {
	Publish(topic string, msg any) error
	PublishBytes(topic string, msg []byte) error
	Subscribe(topic string) (<-chan any, error)
	// SubscribeBytes any is []byte
	SubscribeBytes(topic string) (<-chan any, error)
	Unsubscribe(topic string, sub <-chan any) error
	Close()
	GetPayLoad(sub <-chan any) any
	SetConditions(capacity int)
}


func NewMq() Mq {
	return &MessageQueue[any]{
		closeCh: make(chan struct{}),
	}
	mq.loadFromDisk()

	go mq.startAutoSave() // auto save to disk
	return mq
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
func NewRedisMQ(addr string, password string, db int) (*RedisMq, error) {
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

	return &RedisMq{
		client:   client,
		ctx:      ctx,
		cancel:   cancel,
		channels: make(map[string]chan any),
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
