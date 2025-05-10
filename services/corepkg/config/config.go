package config

import (
	"sync"

	"github.com/aenjoy/iot-lubricant/pkg/types/crypto"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
	"github.com/caarlos0/env/v11"
)

const APP_NAME = "lubricant-core-config"

var _ ioc.Object = (*Config)(nil)

var SystemConfig *Config

type Config struct {
	// app
	AppVersion string
	TlsEnable  bool       `yaml:"tls" env:"TLS_ENABLE" envDefault:"false"`
	HTTPTls    bool       `yaml:"tls_http" env:"HTTP_TLS_ENABLE" envDefault:"false"`
	GRPCTls    bool       `yaml:"tls_grpc" env:"GRPC_TLS_ENABLE" envDefault:"false"`
	Tls        crypto.Tls `yaml:"tls_config" env:"TLS"`

	// grpc
	GrpcPort int `yaml:"grpc_port" env:"GRPC_LISTEN_PORT" envDefault:"5423"`

	// web
	Host    string `yaml:"host" env:"HTTP_LISTEN_HOST" envDefault:"0.0.0.0"`
	WebPort int    `yaml:"web_port" env:"HTTP_LISTEN_PORT" envDefault:"8080"`
	Domain  string `yaml:"domain" env:"HOSTNAME" envDefault:"localhost"`

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
	MqType      string `yaml:"type" env:"MQ_TYPE" envDefault:"internal"` // support: kafka,redis,nats,internal
	KaBrokers   string `yaml:"kafka_brokers" env:"KAFKA_BROKERS"`
	KaGroupID   string `yaml:"kafka_group_id" env:"KAFKA_GROUP_ID"`
	KaPartition int    `yaml:"kafka_partition" env:"KAFKA_PARTITION"`
	NatUrl      string `yaml:"nats_url" env:"NATS_URL"`
	// if set MqType to `redis`,need to set RedisHost,RedisPort,RedisPassword, and RedisDB

	// Auth Provider
	AuthProvider        string `yaml:"auth_provider" env:"AUTH_MODE" envDefault:"casdoor"`
	AuthEndpoint        string `yaml:"auth_endpoint" env:"AUTH_ENDPOINT"`
	AuthClientID        string `yaml:"auth_client_id" env:"AUTH_CLIENT_ID"`
	AuthClientSecret    string `yaml:"auth_client_secret" env:"AUTH_CLIENT_SECRET"`
	AuthOrganization    string `yaml:"auth_organization" env:"AUTH_ORGANIZATION"`
	AuthPublicKeyFile   string `yaml:"auth_public_key_file" env:"AUTH_PUBLICKEYFILE" `
	AuthApplicationName string `yaml:"auth_application_name" env:"AUTH_APPLICATION_NAME" envDefault:"application_lubricant"`

	InternalWorkThreadNumber int    `yaml:"internal_work_thread_number" env:"INTERNAL_WORK_THREAD_NUMBER" envDefault:"4096"`
	EtcdEndpoints            string `yaml:"etcd_endpoints" env:"ETCD_ENDPOINTS"` // "," split
}

func (c *Config) Init() error {
	return env.Parse(c)
}

func (Config) Weight() uint16 {
	return ioc.Config
}

func (c *Config) Version() string {
	return c.AppVersion
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
