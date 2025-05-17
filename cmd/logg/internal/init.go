package internal

import (
	"context"
	"errors"
	"fmt"

	mqV2 "github.com/aenjoy/iot-lubricant/pkg/utils/mq/v2"
	"github.com/aenjoy/iot-lubricant/services/logg/config"
	"github.com/aenjoy/iot-lubricant/services/logg/dao"
)

func GetDb() dao.ILogg {
	return dao.LogDatabase()
}
func GetMq() (mqV2.Mq, error) {
	c := config.GetConfig()
	switch c.MqType {
	case "kafka":
		return nil, errors.New("kafka is unsupported now ")
	case "redis":
		return mqV2.NewRedisMq(context.Background(), mqV2.RedisMqOptions{
			Addrs:    []string{fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort)},
			Password: c.RedisPassword,
			DB:       c.RedisDB,
		})
	case "nats":
		return mqV2.NewNatsMq(c.NatUrl)
	default:
		return nil, errors.New("mq type error")
	}
}
