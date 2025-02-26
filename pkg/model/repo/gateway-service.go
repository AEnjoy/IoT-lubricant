package repo

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	model2 "github.com/AEnjoy/IoT-lubricant/pkg/model"
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
	err = db.AutoMigrate(&model2.Agent{}, &model2.ServerInfo{}, &model2.Gateway{}, &model2.AgentInstance{})
	if err != nil {
		logger.Fatalf("failed to migrate database: %v", err)
	}
	return &GatewayDb{db: db}
}
