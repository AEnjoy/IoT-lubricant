package edge

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

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
	if !a.checkConfigInvalidGet() {
		return ErrInvalidConfig
	}

	paths := a.OpenApi.GetPaths()
	cycle := time.Duration(int64(a.config.Cycle) * int64(time.Second))

	ctx, cancel := context.WithCancel(ctx)
	a.cancel = cancel

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

func (a *app) StopGather(context.Context) error {
	a.cancel()
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
func (a *app) checkConfigInvalidGet() bool {
	// 检查至少一个选项启用且配置有效
	for _, item := range a.GetPaths() {
		opera := item.GetGet()
		if item.GetPost() != nil && opera == nil { // POST
			continue
		}

		if opera == nil {
			return false
		}
		parameters := opera.GetParameters()
		for _, param := range parameters {
			t := param.Schema.GetProperties()[param.Name].Type
			if t == "" && param.Required {
				return false
			}
		}
	}
	return true
}
