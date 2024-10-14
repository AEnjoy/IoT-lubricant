package model

import (
	"gorm.io/gorm"
)

type CoreDb struct {
	db *gorm.DB
}

func (d *CoreDb) IsGatewayIdExists(id string) bool {
	return d.db.Where("id = ?", id).First(&Gateway{}).Error == nil
}

func (*CoreDb) Name() string {
	return "Core-database-client"
}
func (d *CoreDb) Init() error {
	d.db = DefaultCoreClient().db
	return nil
}
