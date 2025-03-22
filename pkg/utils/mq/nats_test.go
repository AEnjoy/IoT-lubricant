//go:build !jetstream

package mq

import (
	"context"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/assert"
)

func natsUrl(t *testing.T) (func(), string) {
	url, ok := os.LookupEnv("TEST_NATS_MQ_URL")
	if !ok {
		t.Log("env TEST_NATS_MQ_URL not found, use internal nats")
		opts := &server.Options{
			Port:    4222 + rand.Intn(100),
			Debug:   true,
			Trace:   true,
			LogFile: "nats.log",
		}
		natsServer, err := server.NewServer(opts)
		assert.NoError(t, err)

		t.Log("starting nats server")
		go natsServer.Start()
		if !natsServer.ReadyForConnections(10 * time.Second) {
			t.Fatal("nats server did not start")
		}

		url = natsServer.ClientURL()
		return func() {
			t.Log("shutting down nats server")
			natsServer.Shutdown()
		}, url
	}
	return func() {}, url
}

func TestNatsMqAbility(t *testing.T) {
	assert := assert.New(t)
	clean, url := natsUrl(t)
	defer clean()

	mq, err := NewNatsMq[string](url)
	assert.NoError(err)
	defer mq.Close()
	mq.SetConditions(100)
	t.Run("Production and re consumption", func(t *testing.T) {
		t.Skip("Only in JetStream Enabled can pass this test")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for i := 0; i < 100; i++ {
			assert.NoError(mq.PublishBytes("/test/topic/123", []byte("hello")))
		}
		time.Sleep(time.Second)
		for i := 0; i < 10; i++ {
			go consumerNatsClient(t, url, i, ctx)
		}
		time.Sleep(3 * time.Second)

		var recv = 0
		subTestResult.Range(func(key, value any) bool {
			assert.Equal(value, struct{}{})
			recv++
			return true
		})
		assert.Equal(10, recv)
	})

	subTestResult = sync.Map{}
	t.Run("Consumption regeneration", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for i := 0; i < 10; i++ {
			go consumerNatsClient(t, url, i, ctx)
		}
		time.Sleep(time.Second)
		for i := 0; i < 100; i++ {
			assert.NoError(mq.PublishBytes("/test/topic/123", []byte("hello")))
		}
		time.Sleep(3 * time.Second)
		var recv = 0
		subTestResult.Range(func(key, value any) bool {
			assert.Equal(value, struct{}{})
			recv++
			return true
		})
		assert.Equal(10, recv)
	})
}
