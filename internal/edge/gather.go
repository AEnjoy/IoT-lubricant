package edge

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
	"github.com/AEnjoy/IoT-lubricant/pkg/edge/config"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
)

var (
	ErrInvalidConfig = errors.New("invalid config. Please check if all necessary settings have been set")
)

var _ gather = (*app)(nil)

type gather interface {
	StartGather(context.Context) error
	StopGather(context.Context) error
	SaveConfig() error
}

func (a *app) StartGather(ctx context.Context) error { // Get
	a.l.Lock()
	defer a.l.Unlock()
	if !edge.CheckConfigInvalidGet(a) {
		return ErrInvalidConfig
	}

	paths := a.OpenApi.GetPaths()
	cycle := time.Duration(int64(a.config.Cycle) * int64(time.Second))

	for {
		select {
		case <-ctx.Done():
			logger.Info("gather worker routine canceled")
			return nil
		case <-time.Tick(cycle):
			logger.Debugln("进行了一次采集")
			for path, item := range paths {
				operation := item.GetGet()
				if operation != nil {
					data, err := a.SendGETMethod(path, operation.GetParameters())
					if err != nil {
						return err
					}
					logger.Debugln(string(data))
					dataSetCh <- data
				}
			}
		}
	}
}
func (a *app) handelGatherSignalCh() error {
	for {
		select {
		case <-a.ctrl.Done():
			return nil
		case ctx := <-config.GatherSignal:
			if err := a.StartGather(ctx); err != nil {
				return err
			}
		case ctx := <-config.StopSignal:
			if err := a.StopGather(ctx); err != nil {
				return err
			}
		}
	}
}
func (a *app) StopGather(context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorln(r)
			return
		}
	}()
	close(dataSetCh)
	return nil
}
func (a *app) SaveConfig() error {
	f, err := os.OpenFile(a.config.FileName+".enable", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	data, err := json.Marshal(a.OpenApi.(*openapi.ApiInfo).OpenAPICli)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}
