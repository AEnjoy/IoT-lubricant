package repo

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"gorm.io/gorm"
)

var _ ILogg = (*Db)(nil)

type Db struct {
	db *gorm.DB
}

func (d *Db) Write(ctx context.Context, log *model.Log) error {
	return d.db.Model(&model.Log{}).WithContext(ctx).Save(log).Error
}
