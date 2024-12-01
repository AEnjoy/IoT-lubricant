// Edge-Agent
package types

import (
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	"github.com/AEnjoy/IoT-lubricant/protobuf/gateway"
)

type EdgeData gateway.DataMessage

type EdgeSystem struct {
	ID          string `yaml:"id"` // API文档ID
	Description string `yaml:"description"`
	Cycle       int    `yaml:"cycle"`        //采集周期 默认5 单位：秒
	ReportCycle int    `yaml:"report-cycle"` //上报周期 默认30 单位：秒
	Algorithm   string `yaml:"algorithm"`    //压缩算法 '-'不压缩 'gzip' 'lz4' 'zstd'
	EnableSlot  []int  `yaml:"enable-slot"`  // 启用的采集Slot

	FileName     string          `yaml:"file-name"` //ApiDoc本地存储文件路径
	Config       openapi.OpenApi `yaml:"-"`         // original api doc
	EnableConfig openapi.OpenApi `yaml:"-"`         // 用于实现功能的配置
}

type DriverData struct {
	AgentId string `gorm:"column:agent_id"`
	Content []byte `gorm:"column:data;type:bytea;not null"` // 由于获取到的数据不一定是string类型,且可能会经过压缩，所以使用bytes存储

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (DriverData) TableName() string {
	return "driver_data"
}
