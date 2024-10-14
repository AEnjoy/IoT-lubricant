package config

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"github.com/caarlos0/env/v11"
)

const APP_NAME = "lubricant-core-config"

var _ ioc.Object = (*Config)(nil)

type Config struct {
	// app
	AppVersion string

	// web
	Host   string `yaml:"host" env:"HTTP_LISTEN_HOST" envDefault:"0.0.0.0"`
	Port   int    `yaml:"port" env:"HTTP_LISTEN_PORT" envDefault:"8080"`
	Domain string `yaml:"domain" env:"BLOG_HOSTNAME" envDefault:"localhost"`

	// mysql
	MySQLHost     string `yaml:"host" env:"DATASOURCE_HOST,required"`
	MySQLPort     int    `yaml:"port" env:"DATASOURCE_PORT,required"`
	MySQLDB       string `yaml:"database" env:"DATASOURCE_DB,required"`
	MySQLUsername string `yaml:"username" env:"DATASOURCE_USERNAME,required"`
	MySQLPassword string `yaml:"password" env:"DATASOURCE_PASSWORD,required"`
	MySQLDebug    bool   `yaml:"debug" env:"DATASOURCE_DEBUG" envDefault:"false"`

	// redis
	RedisEnable   bool   `yaml:"enable" env:"REDIS_ENABLE" envDefault:"false"`
	RedisHost     string `yaml:"host" env:"REDIS_HOST"`
	RedisPort     int    `yaml:"port" env:"REDIS_PORT"`
	RedisPassword string `yaml:"password" env:"REDIS_PASSWORD"`
	RedisDB       int    `yaml:"db" env:"REDIS_DB"`
}

func (c *Config) Init() error {
	return env.Parse(&c)
}

func (Config) Weight() uint16 {
	return ioc.Config
}

func (c *Config) Version() string {
	return c.AppVersion
}

func init() {
	ioc.Controller.Registry(APP_NAME, &Config{})
}
