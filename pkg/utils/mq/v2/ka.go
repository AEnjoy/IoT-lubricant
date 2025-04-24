package mq

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

var ErrNotAllowBrokersEmpty = errors.New("brokers list cannot be empty")
var ErrNeedGroupID = errors.New("GroupID is required for QueueSubscribe but was not provided in NewKaMq")

// AuthComment:
// I'm not test this code because I don't have a Kafka cluster to test it.
// Who would like to help me complete the test of KaMq and put forward pr on GitHub to remove this comment and explain the test results?

// KaMq implements the Mq interface
type KaMq struct {
	brokers  []string
	groupID  string // Optional GroupID for QueueSubscribe
	capacity int    // Maps to ReaderConfig.QueueCapacity

	writer      *kafka.Writer
	adminClient *kafka.Client // For topic creation

	// Manage readers for Subscribe (no GroupID)
	subscribeReaders map[string]*subscribeReaderEntry
	// Manage readers for QueueSubscribe (with GroupID)
	queueSubscribeReaders map[string]*queueSubscribeReaderEntry

	mu sync.Mutex // Protects reader maps
}

// Helper struct to manage a Subscribe reader lifecycle
type subscribeReaderEntry struct {
	reader    *kafka.Reader
	messages  chan []byte
	cancel    context.CancelFunc
	closeOnce sync.Once // Ensures channel and reader are closed only once
}

// Helper struct to manage a QueueSubscribe reader lifecycle
type queueSubscribeReaderEntry struct {
	reader    *kafka.Reader
	messages  chan any // Interface requires chan any
	cancel    context.CancelFunc
	closeOnce sync.Once // Ensures channel and reader are closed only once
}

// NewKaMq creates a new KaMq instance.
// brokers: list of Kafka broker addresses.
// groupID: optional GroupID for QueueSubscribe. Pass "" if not needed.
func NewKaMq(brokers []string, groupID string) (*KaMq, error) {
	if len(brokers) == 0 {
		return nil, ErrNotAllowBrokersEmpty
	}

	// Configure the Kafka Writer
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.Hash{}, // Or &kafka.RoundRobin{}
		RequiredAcks: kafka.RequireAll,
		Completion:   nil, // Optional: Add a completion function for async writes
		// Other producer configurations can be added here, e.g., BatchSize, BatchTimeout
	}

	// Configure the Kafka Admin Client (for topic creation)
	// Note: Admin client might need different connection settings like SASL/TLS
	// depending on your Kafka cluster configuration.
	// This example assumes basic TCP connection.
	adminClient := &kafka.Client{
		Addr: kafka.TCP(brokers...),
		// For SASL:
		// Dial: (&kafka.Dialer{
		// 	Timeout:   10 * time.Second,
		// 	DualStack: true,
		// 	SASLMechanism: plain.Mechanism{Username: "user", Password: "password"},
		// }).DialFunc,
	}

	mq := &KaMq{
		brokers:               brokers,
		groupID:               groupID,
		capacity:              defaultChannelCapacity,
		writer:                writer,
		adminClient:           adminClient,
		subscribeReaders:      make(map[string]*subscribeReaderEntry),
		queueSubscribeReaders: make(map[string]*queueSubscribeReaderEntry),
	}

	return mq, nil
}

// SetConditions sets conditions for the Mq instance.
// Currently, maps capacity to Kafka Reader's QueueCapacity.
// Note: This only affects readers created *after* this call.
func (k *KaMq) SetConditions(capacity int) {
	k.capacity = capacity
}

// ensureTopicExists checks if a topic exists and creates it if not.
// Returns nil if the topic exists or was successfully created.
// Returns an error otherwise.
func (k *KaMq) ensureTopicExists(topic string) error {
	// Check if the topic already exists is complex and error-prone with kafka-go Admin Client.
	// A simpler approach is to attempt creation and ignore the TopicAlreadyExists error.
	// This might involve checking metadata first in a more robust implementation.
	// For simplicity here, we directly attempt creation.

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1, // Default: 1 partition
			ReplicationFactor: 1, // Default: 1 replication factor
			// Add other topic settings if needed
		},
	}

	// Create topic requires a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := k.adminClient.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Addr:   k.adminClient.Addr,
		Topics: topicConfigs,
	})
	if err != nil {
		// Ignore TopicAlreadyExists error
		if errors.Is(err, kafka.TopicAlreadyExists) {
			fmt.Printf("Topic %s already exists, proceeding.", topic)
			return nil
		}
		// Return other errors
		fmt.Printf("Error creating topic %s: %v", topic, err)
		return fmt.Errorf("failed to create topic %s: %w", topic, err)
	}

	fmt.Printf("Successfully created topic %s", topic)
	// It might take a moment for the topic to be ready after creation.
	// A short delay or polling metadata could be added here for robustness.
	time.Sleep(1 * time.Second) // Small delay

	return nil
}

// Publish sends a message to a specific topic.
func (k *KaMq) Publish(topic string, msg []byte) error {
	// Ensure topic exists before publishing
	err := k.ensureTopicExists(topic)
	if err != nil {
		return fmt.Errorf("failed to ensure topic exists before publishing: %w", err)
	}

	message := kafka.Message{
		Topic: topic,
		Value: msg,
	}
	// Use context.Background() for simplicity, consider using a context tied to KaMq lifecycle
	return k.writer.WriteMessages(context.Background(), message)
}

// Subscribe subscribes to a topic without a consumer group.
// Each message will be received by every subscriber on this topic.
func (k *KaMq) Subscribe(topic string) (<-chan []byte, error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	// Ensure topic exists before subscribing
	err := k.ensureTopicExists(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure topic exists before subscribing: %w", err)
	}

	if _, exists := k.subscribeReaders[topic]; exists {
		return nil, fmt.Errorf("already subscribed to topic %s without a group", topic)
	}

	// Create a new reader for this topic (without GroupID)
	readerConfig := kafka.ReaderConfig{
		Brokers: k.brokers,
		Topic:   topic,
		// No GroupID for Subscribe
		MinBytes:      10e3, // 10KB
		MaxBytes:      10e6, // 10MB
		QueueCapacity: k.capacity,
		// Other reader configurations can be added here
	}

	reader := kafka.NewReader(readerConfig)

	// Create channel for messages
	msgChan := make(chan []byte, k.capacity)

	// Create a context to manage the reader's lifecycle
	ctx, cancel := context.WithCancel(context.Background())

	// Store the reader and associated resources
	entry := &subscribeReaderEntry{
		reader:   reader,
		messages: msgChan,
		cancel:   cancel,
	}
	k.subscribeReaders[topic] = entry

	// Start a goroutine to read messages
	go k.runSubscribeReader(ctx, topic, entry)

	return msgChan, nil
}

// runSubscribeReader is a goroutine that reads messages from a non-grouped reader.
func (k *KaMq) runSubscribeReader(ctx context.Context, topic string, entry *subscribeReaderEntry) {
	defer entry.closeOnce.Do(func() {
		// This function runs only once after the reader stops
		close(entry.messages) // Close the channel first
		err := entry.reader.Close()
		if err != nil {
			fmt.Printf("Error closing Subscribe reader for topic %s: %v", topic, err)
		} else {
			fmt.Printf("Subscribe reader closed for topic %s", topic)
		}

		// Clean up the map entry *after* closing the reader and channel
		k.mu.Lock()
		delete(k.subscribeReaders, topic)
		k.mu.Unlock()
	})

	fmt.Printf("Starting Subscribe reader for topic %s", topic)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Subscribe reader context cancelled for topic %s", topic)
			return // Exit the goroutine

		default:
			// Set a read deadline based on context
			// The kafka.Reader.ReadMessage method handles context cancellation internally
			message, err := entry.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, io.EOF) {
					fmt.Printf("Subscribe reader stopped for topic %s due to context cancellation or EOF", topic)
					return // Exit the goroutine on cancellation or EOF
				}
				// Log other errors and continue or implement retry logic
				fmt.Printf("Error reading message from Subscribe reader for topic %s: %v", topic, err)
				// Add a small backoff before retrying read
				time.Sleep(time.Second)
				continue
			}

			// Send message to the channel
			// This will block if the channel buffer is full until a message is consumed
			select {
			case entry.messages <- message.Value:
			case <-ctx.Done():
				fmt.Printf("Subscribe reader stopped for topic %s while trying to send message", topic)
				return // Exit if context is cancelled while waiting to send
			default:
				// This case is technically not needed if the channel is buffered and sends block,
				// but good practice if it were non-blocking or channel was full.
				// Given the channel is buffered, this branch is not easily reachable unless
				// a message is read but the consumer is stuck.
			}
		}
	}
}

// Unsubscribe stops the non-grouped subscriber for the given topic.
func (k *KaMq) Unsubscribe(topic string) error {
	k.mu.Lock()
	entry, exists := k.subscribeReaders[topic]
	k.mu.Unlock() // Release lock before cancelling

	if !exists {
		return fmt.Errorf("not subscribed to topic %s without a group", topic)
	}

	// Cancel the context to signal the reader goroutine to stop
	entry.cancel()

	// The goroutine will handle closing the channel and reader and removing from map

	return nil
}

// QueuePublish sends a message to a specific topic (functionally same as Publish).
func (k *KaMq) QueuePublish(topic string, msg []byte) error {
	// QueuePublish is producer-side, the "Queue" aspect is consumer-side (GroupID).
	// So, it's the same as Publish from the producer's perspective.
	return k.Publish(topic, msg)
}

// QueueSubscribe subscribes to a topic using the instance's GroupID.
// Messages are distributed among instances with the same GroupID.
func (k *KaMq) QueueSubscribe(topic string) (<-chan any, error) {
	if k.groupID == "" {
		return nil, ErrNeedGroupID
	}

	k.mu.Lock()
	defer k.mu.Unlock()

	// Ensure topic exists before subscribing
	err := k.ensureTopicExists(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure topic exists before queue subscribing: %w", err)
	}

	if _, exists := k.queueSubscribeReaders[topic]; exists {
		return nil, fmt.Errorf("already queue subscribed to topic %s with group %s", topic, k.groupID)
	}

	// Create a new reader for this topic (with GroupID)
	readerConfig := kafka.ReaderConfig{
		Brokers: k.brokers,
		GroupID: k.groupID, // Use the instance's GroupID
		Topic:   topic,
		// CommitInterval: 1 * time.Second, // Enable automatic commits (default is often 0, meaning manual)
		MinBytes:      10e3, // 10KB
		MaxBytes:      10e6, // 10MB
		QueueCapacity: k.capacity,
		// Other reader configurations can be added here
	}

	reader := kafka.NewReader(readerConfig)

	// Create channel for messages (must be chan any)
	msgChan := make(chan any, k.capacity)

	// Create a context to manage the reader's lifecycle
	ctx, cancel := context.WithCancel(context.Background())

	// Store the reader and associated resources
	entry := &queueSubscribeReaderEntry{
		reader:   reader,
		messages: msgChan,
		cancel:   cancel,
	}
	k.queueSubscribeReaders[topic] = entry

	// Start a goroutine to read messages
	go k.runQueueSubscribeReader(ctx, topic, entry)

	return msgChan, nil
}

// runQueueSubscribeReader is a goroutine that reads messages from a grouped reader.
func (k *KaMq) runQueueSubscribeReader(ctx context.Context, topic string, entry *queueSubscribeReaderEntry) {
	defer entry.closeOnce.Do(func() {
		// This function runs only once after the reader stops
		close(entry.messages) // Close the channel first
		err := entry.reader.Close()
		if err != nil {
			fmt.Printf("Error closing QueueSubscribe reader for topic %s (group %s): %v", topic, k.groupID, err)
		} else {
			fmt.Printf("QueueSubscribe reader closed for topic %s (group %s)", topic, k.groupID)
		}

		// Clean up the map entry *after* closing the reader and channel
		k.mu.Lock()
		delete(k.queueSubscribeReaders, topic)
		k.mu.Unlock()
	})

	fmt.Printf("Starting QueueSubscribe reader for topic %s (group %s)", topic, k.groupID)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("QueueSubscribe reader context cancelled for topic %s (group %s)", topic, k.groupID)
			return // Exit the goroutine

		default:
			// Set a read deadline based on context
			message, err := entry.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, io.EOF) {
					fmt.Printf("QueueSubscribe reader stopped for topic %s (group %s) due to context cancellation or EOF", topic, k.groupID)
					return // Exit the goroutine on cancellation or EOF
				}
				// Log other errors and continue or implement retry logic
				fmt.Printf("Error reading message from QueueSubscribe reader for topic %s (group %s): %v", topic, k.groupID, err)
				// Add a small backoff before retrying read
				time.Sleep(time.Second)
				continue
			}

			// Send message value (as any) to the channel
			select {
			case entry.messages <- any(message.Value):
				// Message sent successfully
				// kafka-go reader with GroupID auto-commits by default after a message is read
			case <-ctx.Done():
				fmt.Printf("QueueSubscribe reader stopped for topic %s (group %s) while trying to send message", topic, k.groupID)
				return // Exit if context is cancelled while waiting to send
			default:
				// See comment in runSubscribeReader regarding this default case
			}
		}
	}
}

// QueueUnsubscribe stops the grouped subscriber for the given topic.
func (k *KaMq) QueueUnsubscribe(topic string) error {
	if k.groupID == "" {
		return ErrNeedGroupID
	}

	k.mu.Lock()
	entry, exists := k.queueSubscribeReaders[topic]
	k.mu.Unlock() // Release lock before cancelling

	if !exists {
		return fmt.Errorf("not queue subscribed to topic %s with group %s", topic, k.groupID)
	}

	// Cancel the context to signal the reader goroutine to stop
	entry.cancel()

	// The goroutine will handle closing the channel and reader and removing from map

	return nil
}

// Close closes the producer and all active readers.
func (k *KaMq) Close() {
	// Cancel all subscribe readers
	k.mu.Lock()
	for topic, entry := range k.subscribeReaders {
		fmt.Printf("Cancelling Subscribe reader for topic: %s", topic)
		entry.cancel()
		// Don't delete from map here, the goroutine will do it via defer
	}
	// Cancel all queue subscribe readers
	for topic, entry := range k.queueSubscribeReaders {
		fmt.Printf("Cancelling QueueSubscribe reader for topic: %s (group: %s)", topic, k.groupID)
		entry.cancel()
		// Don't delete from map here, the goroutine will do it via defer
	}
	k.mu.Unlock() // Release lock before closing writer/admin client

	// Close the writer
	err := k.writer.Close()
	if err != nil {
		fmt.Printf("Error closing Kafka writer: %v", err)
	} else {
		fmt.Println("Kafka writer closed.")
	}
}
