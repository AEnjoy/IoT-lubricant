package edge

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApp_StartGather(t *testing.T) {
	t.Skip("Pass")

	t.Log("This test will take about 10s to complete")
	WriteConfig()
	assert := assert.New(t)
	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(8*time.Second))
	app := &app{
		config: testConfig,
		ctrl:   ctx,
		cancel: cf,
	}

	go func() {
		assert.NoError(app.StartGather(ctx))
	}()
}
