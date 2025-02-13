package repo

import (
	"fmt"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/config"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	taskTypes "github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Core(database *gorm.DB) *CoreDb {
	if database != nil {
		return &CoreDb{db: database}
	}

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
		Logger: logger.DefualtLog(),
	})
	if err != nil {
		logger.Debugln("dsn:", dsn)
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&taskTypes.Task{}, &model.Clean{}, &model.GatewayHost{}, &model.ErrorLogs{},
		&model.User{}, &model.AsyncJob{}, &model.Gateway{})
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
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
		logger.Fatal("failed to connect database")
	}
	err = db.AutoMigrate(&model.Agent{}, &model.ServerInfo{}, &model.Gateway{})
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}
	return &GatewayDb{db: db}
}
