package edge

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApp_StartGather(t *testing.T) {
	t.Log("This test will take about 10s to complete")
	WriteConfig()
	assert := assert.New(t)
	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(8*time.Second))
	app := &app{
		config:  testConfig,
		ctrl:    ctx,
		cancel:  cf,
		OpenApi: openAPIConfig,
	}

	go func() {
		assert.NoError(app.StartGather(ctx))
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case v := <-dataSetCh:
			str := string(v)
			assert.Contains(str, "string-key1-string-key2")
		}
	}
}
