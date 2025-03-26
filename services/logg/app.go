package logg

import (
	"context"
	"errors"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	"github.com/aenjoy/iot-lubricant/services/logg/dao"
)

type app struct {
	ctx context.Context
	mq  mq.Mq
	db  dao.ILogg
}

func (a *app) Run() error {
	ch, err := a.mq.SubscribeBytes(constant.MESSAGE_SVC_LOGGER)
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
			go a.handel(data.([]byte))
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
func UseMq(mq mq.Mq, err error) func(*app) error {
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
