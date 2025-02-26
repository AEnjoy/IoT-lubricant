package cache

import (
	"context"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/types/errs"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
)

var (
	_ CacheCli[any] = (*NilCache[any])(nil)
	_ ioc.Object    = (*NilCache[any])(nil)
)

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

func (n NilCache[T]) SetEx(context.Context, string, T, time.Duration) error {
	return errs.ErrNullCache
}

func (n NilCache[T]) Set(context.Context, string, T) error {
	return errs.ErrNullCache
}

func (n NilCache[T]) HSet(context.Context, string, string, T) error {
	return errs.ErrNullCache
}

func (n NilCache[T]) HGet(context.Context, string, string) (T, error) {
	var zero T
	return zero, errs.ErrNullCache
}

func (n NilCache[T]) Get(context.Context, string) (T, error) {
	var zero T
	return zero, errs.ErrNullCache
}

func (n NilCache[T]) Incr(context.Context, string) (int64, error) {
	return 0, errs.ErrNullCache
}

func (n NilCache[T]) Decr(context.Context, string) (int64, error) {
	return 0, errs.ErrNullCache
}

func (n NilCache[T]) Delete(context.Context, string) error {
	return errs.ErrNullCache
}

func (n NilCache[T]) Expire(context.Context, string, time.Duration) error {
	return errs.ErrNullCache
}

func (n NilCache[T]) Close(context.Context) error {
	return nil
}
func NewNullCache[T any]() *NilCache[T] {
	return &NilCache[T]{}
}
