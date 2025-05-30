package cache

import (
	"context"
	"time"

	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
)

const APP_NAME = ioc.APP_NAME_CORE_CACHE

var NeverExpired = time.Unix(0, 0)

type CacheCli[T any] interface {
	SetEx(ctx context.Context, key string, value T, duration time.Duration) error
	Set(ctx context.Context, key string, value T) error
	HSet(ctx context.Context, key string, field string, value T) error
	HGet(ctx context.Context, key string, field string) (T, error)
	Get(ctx context.Context, key string) (T, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	Delete(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, duration time.Duration) error
	Close(ctx context.Context) error
}
