package cache

import (
	"context"
	"errors"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/ioc"
)

var (
	_ CacheCli[any] = (*NilCache[any])(nil)
	_ ioc.Object    = (*NilCache)(nil)
)

var ErrNullCache = errors.New("cache client is nil")

type NilCache[T any] struct {
}

func (n NilCache[T]) Init() error {
	return nil
}

func (n NilCache[T]) Weight() uint16 {
	return ioc.CacheCli
}

func (n NilCache[T]) Version() string {
	return "dev"
}

func (n NilCache[T]) SetEx(ctx context.Context, key string, value T, duration time.Duration) error {
	return ErrNullCache
}

func (n NilCache[T]) Set(ctx context.Context, key string, value T) error {
	return ErrNullCache
}

func (n NilCache[T]) HSet(ctx context.Context, key string, field string, value T) error {
	return ErrNullCache
}

func (n NilCache[T]) HGet(ctx context.Context, key string, field string) (T, error) {
	var zero T
	return zero, ErrNullCache
}

func (n NilCache[T]) Get(ctx context.Context, key string) (T, error) {
	var zero T
	return zero, ErrNullCache
}

func (n NilCache[T]) Incr(ctx context.Context, key string) (int64, error) {
	return 0, ErrNullCache
}

func (n NilCache[T]) Decr(ctx context.Context, key string) (int64, error) {
	return 0, ErrNullCache
}

func (n NilCache[T]) Delete(ctx context.Context, key string) error {
	return ErrNullCache
}

func (n NilCache[T]) Expire(ctx context.Context, key string, duration time.Duration) error {
	return ErrNullCache
}

func (n NilCache[T]) Close(ctx context.Context) error {
	return nil
}
func NewNullCache[T any]() *NilCache[T] {
	return &NilCache[T]{}
}
