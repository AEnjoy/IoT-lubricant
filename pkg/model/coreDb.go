package model

import (
	"gorm.io/gorm"
)

type CoreDb struct {
	db *gorm.DB
}

func (*CoreDb) Name() string {
	return "Core-database-client"
}
func (d *CoreDb) Init() error {
	d.db = DefaultCoreClient().db
	return nil
}
