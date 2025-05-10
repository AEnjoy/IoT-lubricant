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

func (Project) TableName() string {
	return "project"
}

type DataStoreEngine struct {
	ID          int    `gorm:"column:id;primary_key;autoIncrement"`
	ProjectID   string `json:"project_id" gorm:"column:project_id;index"`
	Description string `gorm:"column:description" json:"description"`
	Table       string `gorm:"column:table" json:"table"`

	DataBaseType string `gorm:"column:data_base_type;enum('mysql', 'TDEngine', 'mongodb');default:mysql" json:"data_basse_type"`
	DSN          string `gorm:"column:dsn" json:"dsn"` //base64 Encoded 如果是mysql,则是mysql-dsn,如果是tdengine,则是 LinkerInfo

	CreatedAt time.Time    `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeleteAt  sql.NullTime `json:"deleteAt" gorm:"column:deleted_at;type:datetime"`
}

func (DataStoreEngine) TableName() string {
	return "data_store_engine"
}

type LinkerInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Db   string `json:"db"`

	Schemaless *bool `json:"schemaless,omitempty"` // 是否无模式(TDEngine)
}
