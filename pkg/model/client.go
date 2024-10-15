package model

import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/cmd/core/app/config"
	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Core(database *gorm.DB) *CoreDb {
	if database != nil {
		return &CoreDb{db: database}
	}
	c := ioc.Controller.Get(config.APP_NAME).(*config.Config)
	var (
		user         = c.MySQLUsername
		password     = c.MySQLPassword
		address      = c.MySQLHost
		port         = c.MySQLPort
		databaseName = c.MySQLDB
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local", user, password, address, port, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.DefualtLog(),
	})

	if err != nil {
		panic("failed to connect database")
	}

	return &CoreDb{db: db}
}

func DefaultCoreClient() *CoreDb {
	return Core(nil)
}

func NewGatewayDb(database *gorm.DB) *GatewayDb {
	if database != nil {
		return &GatewayDb{db: database}
	}
	db, err := gorm.Open(sqlite.Open("gateway.db"), &gorm.Config{
		Logger: logger.DefualtLog(),
	})
	if err != nil {
		panic("failed to connect database")
	}

	return &GatewayDb{db: db}
}
