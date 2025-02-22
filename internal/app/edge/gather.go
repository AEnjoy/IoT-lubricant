package edge

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
	"github.com/AEnjoy/IoT-lubricant/pkg/edge/config"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/errs"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	json "github.com/bytedance/sonic"
)

var _ gather = (*app)(nil)

type gather interface {
	StartGather(context.Context) error
	StopGather(context.CancelFunc) error
	SaveConfig() error
}

func (a *app) StartGather(ctx context.Context) error { // Get
	if !config.GatherLock.TryLock() {
		return errs.ErrMultGatherInstance
	}
	defer config.GatherLock.Unlock()

	if !edge.CheckConfigInvalid(config.Config.Config) {
		return errs.ErrInvalidConfig
	}

	enable := config.Config.Config.GetEnable()
	cycle := time.Duration(int64(a.config.Cycle) * int64(time.Second))

	ticker := time.NewTicker(cycle)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("gather worker routine canceled")
			return nil
		case <-ticker.C:
			for slot := range enable.Slot {
				method, path := enable.SlotGetEnable(slot)
				logger.Debugln("进行了一次采集")
				switch method {
				case http.MethodGet:
					data, err := a.config.Config.SendGETMethod(path, enable.Get[path])
					if err != nil {
						return err
					}
					logger.Debugln(string(data))
					dataHandlerCh <- &dataHandler{
						slot, &data,
					}
				case http.MethodPost:
					data, err := a.config.Config.SendPOSTMethod(path, *enable.Post[path])
					if err != nil {
						return err
					}
					logger.Debugln(string(data))
					dataHandlerCh <- &dataHandler{
						slot, &data,
					}
				}
			}
		}
	}
}
func (a *app) handelGatherSignalCh() error {
	errCh := make(chan error)
	defer close(errCh)
	var c context.Context
	var cancel context.CancelFunc
	for {
		select {
		case <-a.ctrl.Done():
			return nil
		case ctx := <-config.GatherSignal:
			go func() {
				c, cancel = context.WithCancel(ctx)
				if err := a.StartGather(c); err != nil {
					errCh <- err
				}
			}()
		case <-config.StopSignal:
			go func() {
				if err := a.StopGather(cancel); err != nil {
					errCh <- err
				}
			}()
		case err := <-errCh:
			return err
		}
	}
}
func (a *app) StopGather(cancel context.CancelFunc) error {
	cancel()
	return nil
}
func (a *app) SaveConfig() error {
	f, err := os.OpenFile(a.config.FileName+".enable", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	data, err := json.Marshal(a.config.Config.(*openapi.ApiInfo).OpenAPICli)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}
