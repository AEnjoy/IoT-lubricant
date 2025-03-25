package internal

import (
	"errors"
	"fmt"

	"github.com/aenjoy/iot-lubricant/pkg/utils/mq"
	"github.com/aenjoy/iot-lubricant/services/logg/config"
	"github.com/aenjoy/iot-lubricant/services/logg/repo"
)

func GetDb() repo.ILogg {
	return repo.LogDatabase()
}
func GetMq() (mq.Mq, error) {
	c := config.GetConfig()
	switch c.MqType {
	case "kafka":
		return nil, errors.New("kafka is unsupported now ")
	case "redis":
		return mq.NewRedisMQ(fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort), c.RedisPassword, c.RedisDB)
	case "nats":
		return mq.NewNatsMq(c.NatUrl)
	default:
		return nil, errors.New("mq type error")
	}
}
