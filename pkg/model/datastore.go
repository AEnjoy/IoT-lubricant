package model

import (
	"database/sql"
	"time"
)

type Project struct {
	ID          int    `gorm:"column:id;primary_key;autoIncrement"`
	ProjectID   string `gorm:"column:project_id;type:varchar(36);uniqueIndex"` // xid
	ProjectName string `gorm:"column:project_name;not null" json:"project_name"`
	Description string `gorm:"column:description" json:"description"`
	UserID      string `gorm:"column:user_id;index;not null"`

	CreatedAt time.Time    `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeleteAt  sql.NullTime `json:"deleteAt" gorm:"column:deleted_at;type:datetime"`
}
type DataStoreEngine struct {
	ID          int    `gorm:"column:id;primary_key;autoIncrement"`
	ProjectID   string `json:"project_id" gorm:"column:project_id;index"`
	Description string `gorm:"column:description" json:"description"`

	DataBasseType string `gorm:"column:data_basse_type;enum('mysql', 'TDEngine', 'mongodb');default:mysql" json:"data_basse_type"`
	DSN           string `gorm:"column:dsn" json:"dsn"` //base64 Encoded

	CreatedAt time.Time    `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeleteAt  sql.NullTime `json:"deleteAt" gorm:"column:deleted_at;type:datetime"`
}
