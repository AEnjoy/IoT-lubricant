package config

import (
	"sync"

	"github.com/caarlos0/env/v11"
)

var SystemConfig *Config

type Config struct {
	// app
	AppVersion string

	// mysql
	MySQLHost     string `yaml:"mysql_host" env:"DB_ADDRESS,required"`
	MySQLPort     int    `yaml:"mysql_port" env:"DB_PORT,required"`
	MySQLDB       string `yaml:"mysql_database" env:"DB_NAME,required"`
	MySQLUsername string `yaml:"mysql_username" env:"DB_USER,required"`
	MySQLPassword string `yaml:"mysql_password" env:"DB_PASSWORD,required"`
	MySQLDebug    bool   `yaml:"mysql_debug" env:"DATASOURCE_DEBUG" envDefault:"false"`

	// redis
	RedisEnable   bool   `yaml:"redis" env:"REDIS_ENABLE" envDefault:"false"`
	RedisHost     string `yaml:"redis_host" env:"REDIS_HOST"`
	RedisPort     int    `yaml:"redis_port" env:"REDIS_PORT"`
	RedisPassword string `yaml:"redis_password" env:"REDIS_PASSWORD"`
	RedisDB       int    `yaml:"redis_db" env:"REDIS_DB"`

	// MessageQueue MQ
	MqType      string `yaml:"type" env:"MQ_TYPE" envDefault:"nats"` // support: kafka,redis,nats
	KaBrokers   string `yaml:"kafka_brokers" env:"KAFKA_BROKERS"`
	KaGroupID   string `yaml:"kafka_group_id" env:"KAFKA_GROUP_ID"`
	KaPartition int    `yaml:"kafka_partition" env:"KAFKA_PARTITION"`
	NatUrl      string `yaml:"nats_url" env:"NATS_URL"`
	// if set MqType to `redis`,need to set RedisHost,RedisPort,RedisPassword, and RedisDB
}

func (c *Config) Init() error {
	return env.Parse(c)
}

var _init sync.Once
var _getConfigLock sync.Mutex

func GetConfig() *Config {
	_getConfigLock.Lock()
	defer _getConfigLock.Unlock()

	if SystemConfig == nil {
		_init.Do(func() {
			SystemConfig = &Config{}
			err := SystemConfig.Init()
			if err != nil {
				panic(err)
			}
		})
	}
	return SystemConfig
}
