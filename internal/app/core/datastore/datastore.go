package datastore

import (
	"fmt"
	"strings"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/config"
	"github.com/AEnjoy/IoT-lubricant/internal/cache"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
)

var _ ioc.Object = (*DataStore)(nil)

type DataStore struct {
	CacheEnable bool
	cache.CacheCli[string]
	repo.ICoreDb
	mq.Mq[[]byte]
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
	c := config.GetConfig()
	switch strings.ToLower(c.MqType) {
	case "nats":
		natsMq, err := mq.NewNatsMq[[]byte](c.NatUrl)
		if err != nil {
			return err
		}
		d.Mq = natsMq
	case "redis":
		redisMq, err := mq.NewRedisMQ[[]byte](fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort), c.RedisPassword, c.RedisDB)
		if err != nil {
			return err
		}
		d.Mq = redisMq
	case "kafka":
		d.Mq = mq.NewKafkaMq[[]byte](c.KaBrokers, c.KaGroupID, c.KaPartition, 10)
	case "internal":
		d.Mq = mq.NewGoMq[[]byte]()
	default:
		return fmt.Errorf("mq type error")
	}
	return nil
}

func (DataStore) Weight() uint16 {
	return ioc.DataStore
}

func (DataStore) Version() string {
	return "dev"
}
