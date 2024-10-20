package config

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/caarlos0/env/v11"
)

const APP_NAME = "lubricant-core-config"

var _ ioc.Object = (*Config)(nil)

type Config struct {
	// app
	AppVersion string
	TlsEnable  bool      `yaml:"tls" env:"TLS_ENABLE" envDefault:"false"`
	HTTPTls    bool      `yaml:"tls_http" env:"HTTP_TLS_ENABLE" envDefault:"false"`
	GRPCTls    bool      `yaml:"tls_grpc" env:"GRPC_TLS_ENABLE" envDefault:"false"`
	Tls        types.Tls `yaml:"tls_config" env:"TLS_CONFIG" envPrefix:"TLS_"`

	// grpc
	GrpcPort int `yaml:"port" env:"GRPC_LISTEN_PORT" envDefault:"9090"`

	// web
	Host    string `yaml:"host" env:"HTTP_LISTEN_HOST" envDefault:"0.0.0.0"`
	WebPort int    `yaml:"port" env:"HTTP_LISTEN_PORT" envDefault:"8080"`
	Domain  string `yaml:"domain" env:"HOSTNAME" envDefault:"localhost"`

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
