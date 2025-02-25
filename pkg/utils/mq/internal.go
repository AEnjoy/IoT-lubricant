package mq

import (
	"errors"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
)

type GoMq[T any] struct {
	mu       sync.RWMutex
	topics   map[string][]chan T
	capacity int
	closed   bool
}

func (mq *GoMq[T]) Publish(topic string, msg T) error {
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	if mq.closed {
		return errors.New("message queue is closed")
	}
	for _, ch := range mq.topics[topic] {
		select {
		case ch <- msg:
		default:
			logger.Error("message queue is full, dropping message:", msg)
		}
	}
	return nil
}

func (mq *GoMq[T]) PublishBytes(topic string, msg []byte) error {
	var data T
	anyData, ok := any(msg).(T)
	if !ok {
		return errors.New("invalid type conversion")
	}
	data = anyData
	return mq.Publish(topic, data)
}

func (mq *GoMq[T]) Subscribe(topic string) (<-chan T, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if mq.closed {
		return nil, errors.New("message queue is closed")
	}

	ch := make(chan T, mq.capacity)
	mq.topics[topic] = append(mq.topics[topic], ch)
	return ch, nil
}

func (mq *GoMq[T]) Unsubscribe(topic string, sub <-chan T) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if mq.closed {
		return errors.New("message queue is closed")
	}

	subscribers, ok := mq.topics[topic]
	if !ok {
		return errors.New("topic not found")
	}
	for i, ch := range subscribers {
		if ch == sub {
			mq.topics[topic] = append(subscribers[:i], subscribers[i+1:]...)
			close(ch)
			return nil
		}
	}
	return errors.New("subscriber not found")
}

func (mq *GoMq[T]) Close() {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if mq.closed {
		return
	}

	mq.closed = true
	for _, subs := range mq.topics {
		for _, ch := range subs {
			close(ch)
		}
	}
	mq.topics = nil
}

func (mq *GoMq[T]) GetPayLoad(sub <-chan T) T {
	return <-sub
}

func (mq *GoMq[T]) SetConditions(capacity int) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if capacity > 0 {
		mq.capacity = capacity
	}
}
