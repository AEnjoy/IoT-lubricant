package repo

import (
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

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
	err = db.AutoMigrate(&model.Agent{}, &model.ServerInfo{}, &model.Gateway{},
		&model.AgentInstance{}, &model.Clean{})
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}
	return &GatewayDb{db: db}
}
