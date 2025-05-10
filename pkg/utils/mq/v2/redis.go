package mq

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrRedisAddressEmpty = errors.New("redis addresses cannot be empty")
var ErrClientIsNil = errors.New("client is nil")

// RedisMq implements the Mq interface using go-redis/v9.
type RedisMq struct {
	client redis.UniversalClient
	ctx    context.Context
	cancel context.CancelFunc

	// For Pub/Sub
	pubsubMap sync.Map // map[string]*redis.PubSub - Stores active Pub/Sub subscriptions

	// For Queues (using Redis Lists and BRPOP)
	queueSubscribers sync.Map       // map[string]chan any - Data channels returned to users
	queueStopChans   sync.Map       // map[string]chan struct{} - Channels to signal queue goroutines to stop
	wg               sync.WaitGroup // Waits for queue subscriber goroutines to finish
	queueCapacity    int            // Capacity for queue subscriber channels
}

// RedisMqOptions holds configuration for creating a new RedisMq.
type RedisMqOptions struct {
	// List of Redis server addresses.
	// If one address is provided, a single-node client is used.
	// If multiple addresses are provided, a cluster client is used.
	Addrs    []string
	Password string
	DB       int // Used only for single-node connections

	// Add other redis.Options or redis.ClusterOptions fields as needed
	// e.g., PoolSize, ReadTimeout, etc.
	PoolSize int
}

// NewRedisMqWithUniversalClient creates a new RedisMq instance using a universal client.
func NewRedisMqWithUniversalClient(ctx context.Context, client redis.UniversalClient) (*RedisMq, error) {
	if client == nil {
		return nil, ErrClientIsNil
	}
	// Ping to ensure connection is established
	if err := client.Ping(ctx).Err(); err != nil {
		// Close client if ping fails
		return nil, fmt.Errorf("failed to connect to redis: %V With Close Result:%v", err, client.Close())
	}

	// Create cancellable context for the MQ instance
	ctx, cancel := context.WithCancel(context.Background())

	mq := &RedisMq{
		client:        client,
		ctx:           ctx,
		cancel:        cancel,
		queueCapacity: defaultChannelCapacity, // Default capacity, can be changed by SetConditions
	}

	return mq, nil
}

// NewRedisMq creates a new RedisMq instance.
// It automatically detects whether to use a single-node or cluster client
// based on the number of addresses provided in RedisMqOptions.
func NewRedisMq(ctx context.Context, opts RedisMqOptions) (*RedisMq, error) {
	if len(opts.Addrs) == 0 {
		return nil, ErrRedisAddressEmpty
	}

	var client redis.UniversalClient

	if len(opts.Addrs) == 1 {
		// Single Node Configuration
		if !strings.Contains(opts.Addrs[0], ":") {
			opts.Addrs[0] = fmt.Sprintf("%s:6379", opts.Addrs[0])
		}
		rdbOpts := &redis.Options{
			Addr:     opts.Addrs[0],
			Password: opts.Password,
			DB:       opts.DB,
		}
		if opts.PoolSize > 0 {
			rdbOpts.PoolSize = opts.PoolSize
		}
		// Add other relevant options from opts if needed
		client = redis.NewClient(rdbOpts)
	} else {
		// Cluster Configuration
		clusterOpts := &redis.ClusterOptions{
			Addrs:    opts.Addrs,
			Password: opts.Password,
			// Add other relevant cluster options from opts if needed
			// Example: Route reads to replicas
			// ReadOnly: true,
			// RouteRandomly: true,
		}
		if opts.PoolSize > 0 {
			clusterOpts.PoolSize = opts.PoolSize
		}

		client = redis.NewClusterClient(clusterOpts)
	}

	return NewRedisMqWithUniversalClient(ctx, client)
}

// Publish sends a message using Redis PUBLISH.
func (r *RedisMq) Publish(topic string, msg []byte) error {
	return r.client.Publish(r.ctx, topic, msg).Err()
}

// Subscribe creates a Pub/Sub subscription.
func (r *RedisMq) Subscribe(topic string) (<-chan any, error) {
	// Check if already subscribed to this topic via Pub/Sub
	if _, loaded := r.pubsubMap.Load(topic); loaded {
		return nil, fmt.Errorf("already subscribed to topic %s via Pub/Sub", topic)
	}

	pubsub := r.client.Subscribe(r.ctx, topic)

	// Wait for confirmation that subscription is created before returning channel
	_, err := pubsub.Receive(r.ctx)
	if err != nil {
		pubsub.Close() // Close if initial receive fails
		return nil, fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	r.pubsubMap.Store(topic, pubsub)

	msgChan := make(chan any, r.queueCapacity) // Use capacity setting
	r.wg.Add(1)                                // Add to wait group for graceful shutdown
	go func() {
		defer r.wg.Done()
		defer close(msgChan) // Ensure channel is closed when goroutine exits
		// Ensure pubsub is eventually closed if this goroutine exits unexpectedly
		// Although Unsubscribe or Close should handle it properly.
		// defer func() {
		// 	if ps, ok := r.pubsubMap.Load(topic); ok {
		// 		ps.(*redis.PubSub).Close()
		// 		r.pubsubMap.Delete(topic)
		// 	}
		// }()

		redisCh := pubsub.Channel()

		for {
			select {
			case <-r.ctx.Done(): // Handle global shutdown
				// Attempt to unsubscribe cleanly before exiting
				if ps, ok := r.pubsubMap.Load(topic); ok {
					err := ps.(*redis.PubSub).Unsubscribe(context.Background(), topic)
					if err != nil {
						fmt.Printf("Error during Unsubscribe for topic '%s': %v\n", topic, err)
					} // Use background context for cleanup
					ps.(*redis.PubSub).Close()
					r.pubsubMap.Delete(topic)
				}
				fmt.Printf("Pub/Sub listener for topic %s stopped due to global context cancellation.\n", topic)
				return // Exit goroutine
			case msg, ok := <-redisCh:
				if !ok {
					// Redis channel closed, likely by Unsubscribe or connection issue
					fmt.Printf("Pub/Sub redis channel for topic %s closed.\n", topic)
					// Remove from map if it wasn't removed by Unsubscribe
					r.pubsubMap.Delete(topic)
					return // Exit goroutine
				}
				// Non-blocking send to avoid blocking redis receiver if user channel is full
				select {
				case msgChan <- []byte(msg.Payload):
					// Message sent successfully
				case <-r.ctx.Done():
					// Global context cancelled while trying to send
					fmt.Printf("Pub/Sub listener for topic %s stopping due to global context cancellation during send.\n", topic)
					return
				default:
					fmt.Printf("Warning: Pub/Sub channel for topic %s buffer full. Message dropped.\n", topic)
				}
			}
		}
	}()

	return msgChan, nil
}

// Unsubscribe stops a Pub/Sub subscription.
func (r *RedisMq) Unsubscribe(topic string) error {
	val, loaded := r.pubsubMap.LoadAndDelete(topic)
	if !loaded {
		return fmt.Errorf("not subscribed to topic %s via Pub/Sub", topic)
	}

	pubsub, ok := val.(*redis.PubSub)
	if !ok {
		// This should ideally not happen due to map type assertion
		return fmt.Errorf("internal error: invalid type found in pubsubMap for topic %s", topic)
	}

	// Unsubscribe from the specific topic first
	errUnsub := pubsub.Unsubscribe(r.ctx, topic)
	// Always close the PubSub object to release resources fully
	errClose := pubsub.Close()

	fmt.Printf("Unsubscribed from Pub/Sub topic: %s\n", topic)

	// Report the first error encountered
	if errUnsub != nil {
		return fmt.Errorf("error unsubscribing from topic %s: %w", topic, errUnsub)
	}
	if errClose != nil {
		return fmt.Errorf("error closing pubsub connection for topic %s: %w", topic, errClose)
	}

	return nil
}

// QueuePublish adds a message to the head of a Redis List (LPUSH).
func (r *RedisMq) QueuePublish(topic string, msg []byte) error {
	// Use LPUSH so BRPOP can retrieve messages in FIFO order (relative to a single producer)
	return r.client.LPush(r.ctx, topic, msg).Err()
}

// QueueSubscribe subscribes to a queue using blocking pop (BRPOP).
func (r *RedisMq) QueueSubscribe(topic string) (<-chan any, error) {
	// Check if already subscribed to this queue
	if _, loaded := r.queueStopChans.Load(topic); loaded {
		return nil, fmt.Errorf("already subscribed to queue %s", topic)
	}

	dataChan := make(chan any, r.queueCapacity)
	stopChan := make(chan struct{})

	// Store channels before starting goroutine to prevent race on check
	r.queueSubscribers.Store(topic, dataChan)
	r.queueStopChans.Store(topic, stopChan)

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		defer func() {
			// Cleanup map entries and close data channel when goroutine exits
			r.queueSubscribers.Delete(topic)
			r.queueStopChans.Delete(topic)
			close(dataChan)
			fmt.Printf("Queue listener goroutine for %s finished.\n", topic)
		}()

		fmt.Printf("Starting queue listener for %s...\n", topic)
		for {
			// Use BRPOP with a 0 timeout to block indefinitely until a message arrives
			// or the context is cancelled.
			result, err := r.client.BRPop(r.ctx, 0*time.Second, topic).Result()

			// Check for errors *after* checking stop conditions
			// to prioritize graceful shutdown signals.

			// Check if told to stop or if global context is cancelled
			select {
			case <-stopChan:
				fmt.Printf("Queue listener for %s received stop signal.\n", topic)
				return // Exit goroutine cleanly
			case <-r.ctx.Done():
				fmt.Printf("Queue listener for %s stopping due to global context cancellation.\n", topic)
				return // Exit goroutine cleanly
			default:
				// Continue processing if no stop signal
			}

			// Now handle BRPOP result/error
			if err != nil {
				// If the error is redis.Nil, it might mean BRPOP timed out
				// (shouldn't happen with 0 timeout unless context cancelled)
				// or the context was cancelled. redis-go might return ctx.Err() directly too.
				if errors.Is(err, redis.Nil) || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					// Context likely cancelled, loop will exit on next iteration check
					// Or redis returned Nil unexpectedly, try again.
					// Adding a small sleep might prevent busy-looping if redis.Nil occurs erroneously.
					// time.Sleep(10 * time.Millisecond) // Optional: prevent potential tight loop on redis.Nil
					continue
				}
				// Log other errors and exit the goroutine, closing the user channel
				fmt.Printf("Error reading from queue %s: %v. Stopping listener.\n", topic, err)
				return // Exit goroutine, defer will close channel
			}

			// BRPOP returns ["list-name", "value"]
			if len(result) >= 2 {
				message := []byte(result[1])
				// Send the message to the user channel, but check for stop signals concurrently
				select {
				case dataChan <- message:
					// Message sent successfully
				case <-stopChan:
					fmt.Printf("Queue listener for %s received stop signal while sending. Discarding message.\n", topic)
					return // Exit goroutine
				case <-r.ctx.Done():
					fmt.Printf("Queue listener for %s stopping due to global context cancellation while sending. Discarding message.\n", topic)
					return // Exit goroutine
				}
			} else {
				// Should not happen with successful BRPOP
				fmt.Printf("Warning: Unexpected result from BRPOP for queue %s: %v\n", topic, result)
			}
		}
	}()

	return dataChan, nil
}

// QueueUnsubscribe stops a queue subscriber goroutine.
func (r *RedisMq) QueueUnsubscribe(topic string) error {
	// Load and delete the stop channel first
	val, loaded := r.queueStopChans.LoadAndDelete(topic)
	if !loaded {
		return fmt.Errorf("not subscribed to queue %s", topic)
	}

	// Also remove the data channel entry immediately
	r.queueSubscribers.Delete(topic)

	stopChan, ok := val.(chan struct{})
	if !ok {
		// Should not happen
		return fmt.Errorf("internal error: invalid type found in queueStopChans for topic %s", topic)
	}

	// Close the stop channel to signal the goroutine
	// Use non-blocking close in case it was already closed somehow
	select {
	case <-stopChan:
		// Already closed, do nothing
	default:
		close(stopChan)
	}

	fmt.Printf("Stop signal sent to queue listener for: %s\n", topic)
	// Note: The goroutine might take a moment to fully stop and release its WaitGroup count.
	// Close() will wait for this using r.wg.Wait().

	return nil
}

// SetConditions updates the capacity for new subscriber channels.
func (r *RedisMq) SetConditions(capacity int) {
	if capacity > 0 {
		r.queueCapacity = capacity
	}
}

// Close gracefully shuts down the RedisMq instance.
func (r *RedisMq) Close() {
	fmt.Println("Closing RedisMq...")

	r.cancel()

	r.pubsubMap.Range(func(key, value interface{}) bool {
		topic := key.(string)
		pubsub := value.(*redis.PubSub)
		fmt.Printf("Closing Pub/Sub for topic: %s\n", topic)
		// Ignore errors during shutdown cleanup
		_ = pubsub.Unsubscribe(context.Background()) // Unsubscribe all topics associated with this pubsub
		_ = pubsub.Close()
		r.pubsubMap.Delete(topic) // Ensure removal
		return true               // Continue iterating
	})

	r.queueStopChans.Range(func(key, value interface{}) bool {
		topic := key.(string)
		stopChan := value.(chan struct{})
		// Ensure map entry is removed even if goroutine hasn't exited yet
		r.queueStopChans.Delete(topic)
		r.queueSubscribers.Delete(topic) // Also delete data chan ref
		// Close channel non-blockingly
		select {
		case <-stopChan:
		default:
			close(stopChan)
			fmt.Printf("Sent stop signal to queue listener during Close: %s\n", topic)
		}
		return true // Continue iterating
	})

	fmt.Println("Waiting for background goroutines to stop...")
	r.wg.Wait()
	fmt.Println("All background goroutines stopped.")

	if r.client != nil {
		if err := r.client.Close(); err != nil {
			fmt.Printf("Error closing Redis client: %v\n", err)
		} else {
			fmt.Println("Redis client closed.")
		}
	}
	fmt.Println("RedisMq closed.")
}

// Ensure RedisMq implements Mq interface
var _ Mq = (*RedisMq)(nil)
