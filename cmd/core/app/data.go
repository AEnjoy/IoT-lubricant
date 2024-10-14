package app

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	cache2 "github.com/AEnjoy/IoT-lubricant/pkg/utils/cache"
	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
)

var _ ioc.Object = (*dataStore)(nil)

const maxBuffer = 50

var dataCli = &dataStore{enable: true}

type dataStore struct {
	enable bool
	cache2.CacheCli[string]
	model.CoreDbCli
}

func (d *dataStore) Init() error {
	if dataCli.enable {
		if dataCli.CacheCli == nil {
			o := ioc.Controller.Get(cache2.APP_NAME)
			dataCli.CacheCli = o.(cache2.CacheCli[string])
		}
	} else if dataCli.CacheCli == nil {
		nilCache := cache2.NewNullCache[string]()
		ioc.Controller.Registry(cache2.APP_NAME, nilCache)
		dataCli.CacheCli = nilCache
	}
	return nil
}

func (dataStore) Weight() uint16 {
	//TODO implement me
	panic("implement me")
}

func (dataStore) Version() string {
	//TODO implement me
	panic("implement me")
}

func HandelRecvData(data *core.Data) {

	s := data.String()
	_ = dataCli.HSet(context.Background(), data.GetAgentID(), "latest", s)
}
