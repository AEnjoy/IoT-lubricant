package logg

import (
	"context"
	"errors"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	mqV2 "github.com/aenjoy/iot-lubricant/pkg/utils/mq/v2"
	"github.com/aenjoy/iot-lubricant/services/logg/dao"

	"github.com/panjf2000/ants/v2"
)

type app struct {
	ctx context.Context
	mq  mqV2.Mq
	db  dao.ILogg
}

func (a *app) Run() error {
	pool, err := ants.NewPoolWithFunc(500, a.handel, ants.WithPreAlloc(true))
	if err != nil {
		logger.Fatalf("Failed to create ants pool: %v", err)
	}
	ch, err := a.mq.QueueSubscribe(constant.MESSAGE_SVC_LOGGER)
	if err != nil {
		return err
	}
	for {
		select {
		case <-a.ctx.Done():
			return nil
		case data, ok := <-ch:
			if !ok {
				continue
			}
			_ = pool.Invoke(data)
		}
	}
}
func NewApp(opts ...func(*app) error) *app {
	var app = new(app)
	for _, opt := range opts {
		if err := opt(app); err != nil {
			logger.Fatalf("Failed to apply option: %v", err)
		}
	}
	return app
}
func UseContext(ctx context.Context) func(*app) error {
	return func(app *app) error {
		app.ctx = ctx
		return nil
	}
}
func UseMq(mq mqV2.Mq, err error) func(*app) error {
	return func(app *app) error {
		if err != nil {
			return err
		}
		app.mq = mq
		return nil
	}
}
func UseDb(db dao.ILogg) func(*app) error {
	return func(app *app) error {
		if db == nil {
			return errors.New("failed to connect logger database")
		}
		app.db = db
		return nil
	}
}
