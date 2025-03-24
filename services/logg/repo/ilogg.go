package repo

import (
	"context"

	"github.com/aenjoy/iot-lubricant/pkg/model"
)

type ILogg interface {
	Write(ctx context.Context, log *model.Log) error
}
