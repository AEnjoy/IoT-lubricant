package repo

import (
	"fmt"
	"os"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Core(database *gorm.DB) *CoreDb {
	if database != nil {
		return &CoreDb{db: database}
	}

	var (
		user         = os.Getenv("DB_USER")
		password     = os.Getenv("DB_PASSWORD")
		address      = os.Getenv("DB_ADDRESS")
		port         = os.Getenv("DB_PORT")
		databaseName = os.Getenv("DB_NAME")
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local", user, password, address, port, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.DefualtLog(),
	})

	if err != nil {
		logger.Debugln("dsn:", dsn)
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
		logger.Fatal("failed to connect database")
	}
	err = db.AutoMigrate(&model.Agent{}, &model.ServerInfo{}, &model.Gateway{})
	if err != nil {
		logger.Fatal(err)
	}
	return &GatewayDb{db: db}
}
