package datastore

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/cache"
)

var _ ioc.Object = (*DataStore)(nil)

type DataStore struct {
	CacheEnable bool
	cache.CacheCli[string]
	model.CoreDbCli
}

func (d *DataStore) Init() error {
	if d.CacheEnable {
		if d.CacheCli == nil {
			o := ioc.Controller.Get(cache.APP_NAME)
			d.CacheCli = o.(cache.CacheCli[string])
		}
	} else if d.CacheCli == nil {
		nilCache := cache.NewNullCache[string]()
		ioc.Controller.Registry(cache.APP_NAME, nilCache)
		d.CacheCli = nilCache
	}
	return nil
}

func (DataStore) Weight() uint16 {
	return ioc.CacheCli
}

func (DataStore) Version() string {
	//TODO implement me
	panic("implement me")
}
