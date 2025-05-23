package cache

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/aenjoy/iot-lubricant/services/corepkg/config"
	"github.com/aenjoy/iot-lubricant/services/corepkg/ioc"
	"github.com/redis/go-redis/v9"
)

var (
	_ CacheCli[any] = (*RedisCli[any])(nil)
	_ ioc.Object    = (*RedisCli[any])(nil)
)

type RedisCli[T any] struct {
	rdb *redis.Client

	addr, password string
	port           int
	db             int
	tlsConfig      *tls.Config
}

func (i *RedisCli[T]) Init() error {
	if i.addr == "" || i.password == "" || i.port == 0 {
		conf := config.GetConfig()
		i.addr = conf.RedisHost
		i.port = conf.RedisPort
		i.password = conf.RedisPassword
		i.db = conf.RedisDB
	}

	i.rdb = redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%d", i.addr, i.port),
		Password:  i.password,
		DB:        i.db,
		TLSConfig: i.tlsConfig,
	})
	_, err := i.rdb.Ping(context.Background()).Result()
	return err
}

func (i *RedisCli[T]) Weight() uint16 {
	return ioc.CacheCli
}

func (i *RedisCli[T]) Version() string {
	return "dev"
}

func (r *RedisCli[T]) HSet(ctx context.Context, key string, field string, value T) error {
	return r.rdb.HSet(ctx, key, field, value).Err()
}

func (r *RedisCli[T]) HGet(ctx context.Context, key string, field string) (T, error) {
	var zero T
	val, err := r.rdb.HGet(ctx, key, field).Result()
	if err != nil {
		return zero, err
	}
	return any(val).(T), nil
}

// SetEx 实现设置带过期时间的缓存项
func (r *RedisCli[T]) SetEx(ctx context.Context, key string, value T, duration time.Duration) error {
	return r.rdb.Set(ctx, key, value, duration).Err()
}

// Set 实现设置缓存项
func (r *RedisCli[T]) Set(ctx context.Context, key string, value T) error {
	return r.rdb.Set(ctx, key, value, 0).Err()
}

// Get 实现获取缓存项
func (r *RedisCli[T]) Get(ctx context.Context, key string) (T, error) {
	var zero T // 默认零值处理
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return zero, err
	}

	// 由于 redis 返回的是字符串类型, 需要进一步处理泛型反序列化逻辑
	// 这里假设你的类型可以直接从字符串转换或通过反序列化机制来获取
	// 比如，你可能会使用 encoding/json 来反序列化。
	return any(val).(T), nil
}

// Incr 实现递增缓存项
func (r *RedisCli[T]) Incr(ctx context.Context, key string) (int64, error) {
	return r.rdb.Incr(ctx, key).Result()
}

// Decr 实现递减缓存项
func (r *RedisCli[T]) Decr(ctx context.Context, key string) (int64, error) {
	return r.rdb.Decr(ctx, key).Result()
}

// Delete 实现删除缓存项
func (r *RedisCli[T]) Delete(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

// Expire 实现设置缓存项的过期时间
func (r *RedisCli[T]) Expire(ctx context.Context, key string, duration time.Duration) error {
	return r.rdb.Expire(ctx, key, duration).Err()
}

// Close 实现关闭缓存连接
func (r *RedisCli[T]) Close(ctx context.Context) error {
	return r.rdb.Close()
}

// NewRedisCli 创建新的 Redis 客户端实例
func NewRedisCli[T any](addr, password string, db int, tlsConfig *tls.Config) *RedisCli[T] {
	return &RedisCli[T]{addr: addr, password: password, db: db, tlsConfig: tlsConfig}
}
