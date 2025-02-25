package mq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testBenchmark[T any](b *testing.B, mq *KafkaMq[T]) {
	assert := assert.New(b)
	//resultCh := make(chan int, runtime.GOMAXPROCS(0))
	//defer close(resultCh)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			assert.NoError(mq.PublishBytes("test", []byte{2}))
		}
	})
	b.StopTimer()
	subscribe, err := mq.Subscribe("test")
	assert.NoError(err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			assert.NotNil(mq.GetPayLoad(subscribe))
		}
	})
	b.StopTimer()
}
