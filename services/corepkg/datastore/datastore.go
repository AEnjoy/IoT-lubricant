package datastore

import (
	"context"
	"fmt"
	"strings"

	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	mqV2 "github.com/aenjoy/iot-lubricant/pkg/utils/mq/v2"

	"github.com/aenjoy/iot-lubricant/services/corepkg/cache"
	"github.com/aenjoy/iot-lubricant/services/corepkg/config"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
	"github.com/aenjoy/iot-lubricant/services/corepkg/repo"
)

var _ ioc.Object = (*DataStore)(nil)

type DataStore struct {
	CacheEnable bool
	cache.CacheCli[string]
	repo.ICoreDb
	mq.Mq
	V2mq mqV2.Mq
}

func (d *DataStore) Init() error {
	d.ICoreDb = ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(repo.ICoreDb)
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
		natsMq, err := mq.NewNatsMq(c.NatUrl)
		if err != nil {
			return err
		}
		d.Mq = natsMq
		v2, err := mqV2.NewNatsMq(c.NatUrl)
		if err != nil {
			return err
		}
		d.V2mq = v2
	case "redis":
		redisMq, err := mq.NewRedisMQ(fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort), c.RedisPassword, c.RedisDB)
		if err != nil {
			return err
		}
		d.Mq = redisMq
		address := strings.Split(c.RedisHost, ",")
		for i := range address {
			if !strings.Contains(address[i], ":") {
				address[i] = fmt.Sprintf("%s:%d", address[i], c.RedisPort)
			}
		}
		redis, err := mqV2.NewRedisMq(context.Background(), mqV2.RedisMqOptions{
			Addrs:    address,
			Password: c.RedisPassword,
			DB:       c.RedisDB,
		})
		if err != nil {
			return err
		}
		d.V2mq = redis
	//case "kafka":
	//	d.Mq = mq.NewKafkaMq(c.KaBrokers, c.KaGroupID, c.KaPartition, 10)
	//case "internal":
	//	d.Mq = mq.NewGoMq[[]byte]()
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
