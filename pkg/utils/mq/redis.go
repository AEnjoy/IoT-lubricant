package mq

import (
	"context"
	"errors"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

type RedisMq struct {
	client   *redis.Client
	ctx      context.Context
	cancel   context.CancelFunc
	channels map[string]chan any
	subs     map[string]*redis.PubSub
	mutex    sync.RWMutex
	capacity int
}

func (r *RedisMq) Publish(topic string, msg any) error {
	marshal, err := sonic.Marshal(msg)
	if err != nil {
		return err
	}
	return r.PublishBytes(topic, marshal)
}

func (r *RedisMq) Subscribe(topic string) (<-chan any, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.channels[topic]; exists {
		return nil, errors.New("already subscribed to topic")
	}

	sub := r.client.Subscribe(r.ctx, topic)
	ch := make(chan any, r.capacity)

	r.channels[topic] = ch
	r.subs[topic] = sub

	go func() {
		for msg := range sub.Channel() {
			var data any
			if err := sonic.Unmarshal([]byte(msg.Payload), &data); err == nil {
				ch <- data
			}
		}
	}()

	return ch, nil
}

func (r *RedisMq) SubscribeBytes(topic string) (<-chan any, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, exists := r.channels[topic]; exists {
		return nil, errors.New("already subscribed to topic")
	}

	sub := r.client.Subscribe(r.ctx, topic)
	ch := make(chan any, r.capacity)

	r.channels[topic] = ch
	r.subs[topic] = sub

	go func() {
		for msg := range sub.Channel() {
			ch <- msg.Payload
		}
	}()
	return ch, nil
}

func (r *RedisMq) PublishBytes(topic string, msg []byte) error {
	return r.client.Publish(r.ctx, topic, msg).Err()
}

func (r *RedisMq) Unsubscribe(topic string, sub <-chan any) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	ch, exists := r.channels[topic]
	ps, subExists := r.subs[topic]

	if !exists || !subExists || ch != sub {
		return errors.New("subscription not found")
	}

	_ = ps.Close()
	delete(r.subs, topic)

	close(ch)
	delete(r.channels, topic)

	return nil
}

func (r *RedisMq) Close() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.cancel()
	r.client.Close()

	for topic, ch := range r.channels {
		close(ch)
		delete(r.channels, topic)
	}
}

func (r *RedisMq) GetPayLoad(sub <-chan any) any {
	return <-sub
}

func (r *RedisMq) SetConditions(capacity int) {
	r.capacity = capacity
}
