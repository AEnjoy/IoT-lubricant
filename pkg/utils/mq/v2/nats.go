package mq

import (
	"errors"
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
)

var (
	ErrNATSNotActive = errors.New("NATS connection is not active")
)

// NatsMq implements the Mq interface using NATS.
type NatsMq struct {
	conn *nats.Conn
	mu   sync.Mutex

	// Store regular subscriptions (topic -> *nats.Subscription)
	subscriptions map[string]*nats.Subscription
	// Store channels for regular subscriptions (topic -> chan any)
	subChannels map[string]chan any

	// Store queue subscriptions (topic -> *nats.Subscription)
	// Note: NATS uses topic+queueGroup for uniqueness, we map by topic for simplicity
	queueSubscriptions map[string]*nats.Subscription
	// Store channels for queue subscriptions (topic -> chan any)
	queueSubChannels map[string]chan any

	// Capacity for the buffered channels returned to the user
	channelCapacity int
}

// NewNatsMq creates a new NatsMq instance and connects to the NATS server.
// natsURL should be the connection URL (e.g., "nats://localhost:4222").
func NewNatsMq(natsURL string) (*NatsMq, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS at %s: %v", natsURL, err)
	}

	mq := &NatsMq{
		conn:               nc,
		subscriptions:      make(map[string]*nats.Subscription),
		subChannels:        make(map[string]chan any),
		queueSubscriptions: make(map[string]*nats.Subscription),
		queueSubChannels:   make(map[string]chan any),
		channelCapacity:    defaultChannelCapacity, // Set default capacity
	}

	return mq, nil
}

// Publish sends a message to a specific topic using NATS.
func (n *NatsMq) Publish(topic string, msg []byte) error {
	if n.conn == nil || !n.conn.IsConnected() {
		return ErrNATSNotActive
	}

	err := n.conn.Publish(topic, msg)
	if err != nil {
		return fmt.Errorf("NATS publish to topic '%s' failed: %w", topic, err)
	}
	// It's often good practice to Flush after important publishes
	// if you need higher guarantee of delivery attempt before returning.
	// err = n.conn.Flush()
	// if err != nil {
	// 	 return fmt.Errorf("NATS flush after publish to topic '%s' failed: %w", topic, err)
	// }
	return nil // Publish itself doesn't guarantee delivery, Flush does more
}

// Subscribe creates a NATS subscription and returns a channel for messages.
func (n *NatsMq) Subscribe(topic string) (<-chan any, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.conn == nil || !n.conn.IsConnected() {
		return nil, ErrNATSNotActive
	}

	if _, exists := n.subscriptions[topic]; exists {
		return nil, fmt.Errorf("already subscribed to topic '%s'", topic)
	}

	userChan := make(chan any, n.channelCapacity)
	n.subChannels[topic] = userChan // Store before subscribe in case of immediate messages

	// Use ChanSubscribe for easy integration with Go channels.
	// NATS library handles buffering between the network and this channel.
	natsChan := make(chan *nats.Msg, n.channelCapacity*2) // Internal NATS chan, often larger buffer is good
	sub, err := n.conn.ChanSubscribe(topic, natsChan)
	if err != nil {
		// Clean up if subscribe failed
		close(userChan)
		delete(n.subChannels, topic)
		return nil, fmt.Errorf("NATS subscribe to topic '%s' failed: %w", topic, err)
	}

	n.subscriptions[topic] = sub // Store the NATS subscription

	// Start a goroutine to forward messages from natsChan to userChan
	go func() {
		// This loop exits when natsChan is closed (usually by Unsubscribe)
		for msg := range natsChan {
			// Select with a default to prevent blocking forever if userChan buffer is full
			// although with a reasonable buffer size this is less likely.
			// A more robust solution might involve dropping messages or logging if blocked.
			select {
			case userChan <- msg.Data:
			default:
				// Handle buffer full case - log, drop, etc.
				fmt.Printf("Warning: Channel buffer full for topic '%s', message dropped.\n", topic)
			}
		}

		// Once natsChan is closed by NATS (e.g., via sub.Unsubscribe()), close the user channel.
		n.mu.Lock()
		close(userChan)
		delete(n.subChannels, topic)
		n.mu.Unlock()
	}()

	return userChan, nil
}

// Unsubscribe removes a NATS subscription.
func (n *NatsMq) Unsubscribe(topic string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	sub, exists := n.subscriptions[topic]
	if !exists {
		return nil
	}

	err := sub.Unsubscribe()
	if err != nil {
		return fmt.Errorf("NATS unsubscribe for topic '%s' failed: %w", topic, err)
	}

	delete(n.subscriptions, topic)
	// ** Don't close userChan here, the forwarding goroutine does that when natsChan closes. **
	// Just remove it from the map.
	delete(n.subChannels, topic)

	return nil
}

// QueuePublish publishes a message for queue subscribers. In NATS, this is the same as Publish.
func (n *NatsMq) QueuePublish(topic string, msg []byte) error {
	// The "queue" aspect is handled by QueueSubscribe, not the publish call itself.
	return n.Publish(topic, msg)
}

// getQueueGroupName generates a queue group name based on the topic.
// You might want a more sophisticated strategy depending on your application needs.
func (n *NatsMq) getQueueGroupName(topic string) string {
	return topic + defaultQueueGroupSuffix
}

// QueueSubscribe creates a NATS queue subscription.
func (n *NatsMq) QueueSubscribe(topic string) (<-chan any, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.conn == nil || !n.conn.IsConnected() {
		return nil, ErrNATSNotActive
	}

	// Check if a queue subscription for this *topic* already exists in our map
	// Note: NATS allows multiple QueueSubscribe calls with the same topic/group,
	// they just join the group. We prevent multiple calls *from this instance*
	// returning different channels for the same logical topic subscription.
	if _, exists := n.queueSubscriptions[topic]; exists {
		return nil, fmt.Errorf("already queue subscribed to topic '%s'", topic)
	}

	queueGroup := n.getQueueGroupName(topic) // Determine queue group name

	// Create the buffered channel for the user (type 'any')
	userChan := make(chan any, n.channelCapacity)
	n.queueSubChannels[topic] = userChan // Store before subscribe

	// Use ChanQueueSubscribe
	natsChan := make(chan *nats.Msg, n.channelCapacity*2)
	sub, err := n.conn.ChanQueueSubscribe(topic, queueGroup, natsChan)
	if err != nil {
		close(userChan) // Clean up
		delete(n.queueSubChannels, topic)
		return nil, fmt.Errorf("NATS queue subscribe to topic '%s' (group '%s') failed: %w", topic, queueGroup, err)
	}

	n.queueSubscriptions[topic] = sub // Store the NATS subscription

	// Start forwarding goroutine
	go func() {
		for msg := range natsChan {
			// Interface requires 'any'. We send the raw byte slice as 'any'.
			// The consumer will need to type assert if they expect []byte.
			select {
			case userChan <- msg.Data:
			default:
				fmt.Printf("Warning: Queue channel buffer full for topic '%s', message dropped.\n", topic)
			}
		}
		// Once natsChan closes, close userChan
		n.mu.Lock()
		close(userChan)
		delete(n.queueSubChannels, topic) // Clean up map entry
		n.mu.Unlock()
	}()

	return userChan, nil
}

// QueueUnsubscribe removes a NATS queue subscription.
func (n *NatsMq) QueueUnsubscribe(topic string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	sub, exists := n.queueSubscriptions[topic]
	if !exists {
		return nil // Idempotent: Not subscribed
	}

	// NATS Unsubscribe handles closing the associated channel
	err := sub.Unsubscribe()
	if err != nil {
		return fmt.Errorf("NATS queue unsubscribe for topic '%s' failed: %w", topic, err)
	}

	// Clean up tracking maps
	delete(n.queueSubscriptions, topic)
	delete(n.queueSubChannels, topic) // Goroutine will close the actual channel
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	close(n.queueSubChannels[topic])
	return nil
}

// Close gracefully closes the NATS connection and cleans up resources.
func (n *NatsMq) Close() {
	n.mu.Lock() // Lock to prevent concurrent modifications/access
	defer n.mu.Unlock()

	if n.conn == nil {
		return // Already closed or never initialized
	}

	// Unsubscribe from all regular subscriptions
	// Iterate safely over map keys
	topics := make([]string, 0, len(n.subscriptions))
	for topic := range n.subscriptions {
		topics = append(topics, topic)
	}
	n.mu.Unlock() // Unlock temporarily to allow Unsubscribe to lock
	for _, topic := range topics {
		if err := n.Unsubscribe(topic); err != nil {
			fmt.Printf("Error during Close->Unsubscribe for topic '%s': %v\n", topic, err)
		}
	}
	n.mu.Lock() // Re-acquire lock

	// Unsubscribe from all queue subscriptions
	queueTopics := make([]string, 0, len(n.queueSubscriptions))
	for topic := range n.queueSubscriptions {
		queueTopics = append(queueTopics, topic)
	}
	n.mu.Unlock() // Unlock temporarily
	for _, topic := range queueTopics {
		if err := n.QueueUnsubscribe(topic); err != nil {
			fmt.Printf("Error during Close->QueueUnsubscribe for topic '%s': %v\n", topic, err)
		}
	}
	n.mu.Lock() // Re-acquire lock

	n.subscriptions = nil
	n.subChannels = nil
	n.queueSubscriptions = nil
	n.queueSubChannels = nil

	// Close the NATS connection
	if n.conn != nil && !n.conn.IsClosed() {
		n.conn.Close()
		n.conn = nil
	}
}

// SetConditions configures the capacity for newly created subscription channels.
// Note: This does not affect already existing channels.
func (n *NatsMq) SetConditions(capacity int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if capacity <= 0 {
		fmt.Printf("Warning: Invalid channel capacity %d provided, using default %d\n", capacity, defaultChannelCapacity)
		n.channelCapacity = defaultChannelCapacity
	} else {
		n.channelCapacity = capacity
	}
}

// Compile-time check to ensure NatsMq implements Mq
var _ Mq = (*NatsMq)(nil)
