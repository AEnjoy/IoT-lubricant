package repo

import (
	"fmt"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	taskTypes "github.com/aenjoy/iot-lubricant/pkg/types/task"
	"github.com/aenjoy/iot-lubricant/services/lubricant/config"
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
		&model.User{}, &model.AsyncJob{}, &model.Gateway{}, &model.Token{}, &model.Agent{}, &model.Data{},
		&model.GatherNodeConfig{})
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}

	if !db.Migrator().HasIndex(&model.Gateway{}, "idx_user_gateway") {
		err = db.Exec(`
			CREATE UNIQUE INDEX idx_user_gateway ON gateway(user_id, gateway_id);
		`).Error
		if err != nil {
			logger.Fatalf("failed to create unique index: %v", err)
		}
	} else {
		logger.Debugln("Index idx_user_gateway already exists.")
	}
	db.Exec(`UPDATE async_job
				SET status = 'failed',updated_at = NOW()
				WHERE expired_at < NOW() AND status != 'completed';`)

	return &CoreDb{db: db}
}
func DefaultCoreClient() *CoreDb {
	return Core(nil)
}
