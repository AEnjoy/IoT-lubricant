package repo

import (
	"fmt"
	"os"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/services/logg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func LogDatabase() *Db {
	conf := config.GetConfig()
	var (
		user         = conf.MySQLUsername
		password     = conf.MySQLPassword
		address      = conf.MySQLHost
		port         = conf.MySQLPort
		databaseName = conf.MySQLDB
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local", user, password, address, port, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: func() gormLogger.Interface {
			if os.Getenv(constant.ENV_RUNNING_LEVEL) == "debug" {
				return logger.DefualtLog()
			}
			return nil
		}(),
	})
	if err != nil {
		logger.Fatalf("failed to connect database: DSN:%s", dsn)
	}
	err = db.AutoMigrate(&model.Log{})
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}

	return &Db{db: db}
}
