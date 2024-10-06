package mq

import (
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils"
)

var _ Mq[any] = (*MessageQueue[any])(nil)

type MessageQueue[T any] struct {
	topics      sync.Map       // 存储topic对应的订阅者列表
	topicClosed sync.Map       // 标记每个topic是否已关闭
	capacity    int            // 每个channel的容量
	closeCh     chan struct{}  // 用于关闭队列
	wg          sync.WaitGroup // 等待所有goroutine结束
	closeOnce   sync.Once      // 确保队列只关闭一次
}

// SetConditions 设置队列的容量
func (mq *MessageQueue[T]) SetConditions(capacity int) {
	mq.capacity = capacity
}

func (mq *MessageQueue[T]) Publish(topic string, msg T) error {
	value, ok := mq.topics.Load(topic)
	if !ok {
		return ErrNotFoundTopic
	}

	subscribers := value.([]chan T)

	mq.wg.Add(len(subscribers))
	for _, sub := range subscribers {
		go func(s chan T) {
			defer mq.wg.Done()
			if utils.IsClosed(s) {
				return
			}
			select {
			case s <- msg:
			case <-mq.closeCh:
				return
			}
		}(sub)
	}

	return nil
}
func (mq *MessageQueue[T]) PublishBytes(topic string, msg []byte) error {
	var msgT T
	bytesToMsg[T](msg, &msgT)
	return mq.Publish(topic, msgT)
}
func (mq *MessageQueue[T]) Subscribe(topic string) (<-chan T, error) {
	ch := make(chan T, mq.capacity)

	value, _ := mq.topics.LoadOrStore(topic, []chan T{ch})

	subscribers := value.([]chan T)
	mq.topics.Store(topic, append(subscribers, ch))

	return ch, nil
}

func (mq *MessageQueue[T]) Unsubscribe(topic string, sub <-chan T) error {
	value, ok := mq.topics.Load(topic)
	if !ok {
		return ErrNotFoundTopic
	}

	if closed, _ := mq.topicClosed.LoadOrStore(topic, false); closed.(bool) {
		return nil // topic已关闭，跳过处理
	}

	subscribers := value.([]chan T)
	for i, s := range subscribers {
		if s == sub {
			close(s)
			mq.topics.Store(topic, append(subscribers[:i], subscribers[i+1:]...))
			return nil
		}
	}

	return ErrNotFoundSubscriber
}

// GetPayLoad 从订阅者通道中获取消息
func (mq *MessageQueue[T]) GetPayLoad(sub <-chan T) T {
	select {
	case msg := <-sub:
		return msg
	case <-mq.closeCh:
		var zeroValue T
		return zeroValue
	}
}

func (mq *MessageQueue[T]) Close() {
	mq.closeOnce.Do(func() {
		close(mq.closeCh)
		mq.wg.Wait()

		mq.topics.Range(func(key, value interface{}) bool {
			topic := key.(string)

			// 如果topic已经被标记为关闭，跳过处理
			if closed, _ := mq.topicClosed.LoadOrStore(topic, false); closed.(bool) {
				return true
			}

			subscribers := value.([]chan T)

			// 标记topic为关闭状态
			mq.topicClosed.Store(topic, true)

			for _, sub := range subscribers {
				utils.CloseChannel(sub) // 这里的CloseChannel需要安全关闭
			}
			return true
		})
	})
}
