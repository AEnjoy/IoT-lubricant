package mq

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/bytedance/sonic"
	"github.com/segmentio/kafka-go"
)

type KafkaMq[T any] struct {
	address     []string
	groupID     string
	partition   int
	adminClient *kafka.Client
	writers     map[string]*kafka.Writer
	writersMu   sync.Mutex
	subscribers map[string][]*subscriber[T]
	subsMu      sync.Mutex
	timeout     time.Duration
	capacity    int
}

type subscriber[T any] struct {
	reader   *kafka.Reader
	msgChan  chan T
	stopChan chan struct{}
	cancel   context.CancelFunc
}

func (k *KafkaMq[T]) Publish(topic string, msg T) error {
	data, err := sonic.Marshal(msg)
	if err != nil {
		return err
	}
	return k.PublishBytes(topic, data)
}

func (k *KafkaMq[T]) PublishBytes(topic string, msg []byte) error {
	if err := k.ensureTopicExists(topic); err != nil {
		return err
	}

	writer, err := k.getWriter(topic)
	if err != nil {
		return err
	}

	return writer.WriteMessages(context.Background(), kafka.Message{
		Value: msg,
	})
}

func (k *KafkaMq[T]) Subscribe(topic string) (<-chan T, error) {
	ctx, cancel := context.WithCancel(context.Background())
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   k.address,
		GroupID:   k.groupID,
		Topic:     topic,
		MaxWait:   k.timeout,
		Partition: 0,
	})

	msgChan := make(chan T, k.capacity)
	stopChan := make(chan struct{})

	sub := &subscriber[T]{
		reader:   reader,
		msgChan:  msgChan,
		stopChan: stopChan,
		cancel:   cancel,
	}

	k.subsMu.Lock()
	k.subscribers[topic] = append(k.subscribers[topic], sub)
	k.subsMu.Unlock()

	go k.consumeMessages(ctx, sub)

	return msgChan, nil
}

func (k *KafkaMq[T]) Unsubscribe(topic string, sub <-chan T) error {
	k.subsMu.Lock()
	defer k.subsMu.Unlock()

	subs, exists := k.subscribers[topic]
	if !exists {
		return errors.New("topic not found")
	}

	for i, s := range subs {
		if s.msgChan == sub {
			close(s.stopChan)
			s.cancel()
			logger.Info("subscriber closed")
			s.reader.Close()
			logger.Info("subscriber removed")
			k.subscribers[topic] = append(subs[:i], subs[i+1:]...)
			return nil
		}
	}

	return errors.New("subscriber not found")
}

func (k *KafkaMq[T]) Close() {
	k.writersMu.Lock()
	for _, w := range k.writers {
		w.Close()
	}
	k.writers = make(map[string]*kafka.Writer)
	k.writersMu.Unlock()

	k.subsMu.Lock()
	for _, subs := range k.subscribers {
		for _, s := range subs {
			close(s.stopChan)
			s.reader.Close()
		}
	}
	k.subscribers = make(map[string][]*subscriber[T])
	k.subsMu.Unlock()
}

func (k *KafkaMq[T]) GetPayLoad(sub <-chan T) T {
	return <-sub
}

func (k *KafkaMq[T]) SetConditions(capacity int) {
	if capacity <= 0 {
		capacity = 100
	}
	k.timeout = time.Duration(capacity) * time.Second
}

// Helper methods
func (k *KafkaMq[T]) ensureTopicExists(topic string) error {
	ctx, cancel := context.WithTimeout(context.Background(), k.timeout)
	defer cancel()

	_, err := k.adminClient.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{{
			Topic:             topic,
			NumPartitions:     k.partition,
			ReplicationFactor: 1,
		}},
	})

	if err != nil && !isTopicExistsError(err) {
		return err
	}
	return nil
}

func isTopicExistsError(err error) bool {
	var topicError kafka.Error
	if errors.As(err, &topicError) && topicError == kafka.TopicAlreadyExists {
		return true
	}
	return false
}

func (k *KafkaMq[T]) getWriter(topic string) (*kafka.Writer, error) {
	k.writersMu.Lock()
	defer k.writersMu.Unlock()

	if writer, exists := k.writers[topic]; exists {
		return writer, nil
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(k.address...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		WriteTimeout: k.timeout,
	}

	k.writers[topic] = writer
	return writer, nil
}

func (k *KafkaMq[T]) consumeMessages(ctx context.Context, sub *subscriber[T]) {
	defer close(sub.msgChan)

	for {
		select {
		case <-sub.stopChan:
			return
		default:
			m, err := sub.reader.ReadMessage(ctx)
			if err != nil {
				if strings.Contains(err.Error(), "context canceled") {
					logger.Debugf("context canceled")
					return
				}
				logger.Error("error reading message:", err)
				continue
			}

			var msg T
			if err := sonic.Unmarshal(m.Value, &msg); err != nil {
				continue
			}

			sub.msgChan <- msg
		}
	}
}
