package mq

// Mq defines the interface for a message queue.
type Mq interface {
	// Publish sends a message to a specific topic.
	// All active subscribers on that topic will receive the message.
	Publish(topic string, msg []byte) error

	// Subscribe creates a subscription to a topic.
	// It returns a read-only channel where received messages (as 'any', typically []byte) will be sent.
	// Each subscriber instance receives a copy of the message.
	Subscribe(topic string) (<-chan any, error)

	// Unsubscribe removes the subscription for the given topic.
	// The channel returned by Subscribe will be closed.
	Unsubscribe(topic string) error

	// QueuePublish sends a message to a topic associated with a queue group.
	// Functionally often the same as Publish on the publisher side.
	QueuePublish(topic string, msg []byte) error

	// QueueSubscribe creates a subscription to a topic within a queue group.
	// Only one subscriber within the same queue group will receive a given message.
	// It returns a read-only channel where received messages (as 'any', typically []byte) will be sent.
	// The specific queue group name might be derived from the topic or configured internally.
	QueueSubscribe(topic string) (<-chan any, error)

	// QueueUnsubscribe removes the queue subscription for the given topic.
	// The channel returned by QueueSubscribe will be closed.
	QueueUnsubscribe(topic string) error

	// Close cleans up all resources, unsubscribes from all topics, and closes the connection.
	Close()

	// SetConditions allows configuring parameters like channel buffer capacity.
	SetConditions(capacity int)
}

const (
	// Default capacity for the channels returned by Subscribe/QueueSubscribe
	defaultChannelCapacity = 2048
	// Default Queue Group Name Suffix (can be customized)
	// Using a fixed suffix helps differentiate queue groups if needed
	defaultQueueGroupSuffix = "-qgroup"
)
