package cache

import (
	"sync"
	"time"

	def "github.com/aenjoy/iot-lubricant/pkg/default"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
)

// MemoryCache todo:适配 CacheCli
type MemoryCache[T any] struct {
	cacheMap sync.Map //id-*Result[T]
}

func (m *MemoryCache[T]) cleanExpired() {
	m.cacheMap.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(*Result[T])
		if v.expiredAt.Before(time.Now()) && !v.expiredAt.IsZero() {
			logger.Debugf("clean expired cache: %v", k)
			m.cacheMap.Delete(k)
		}
		return true
	})
}
func (m *MemoryCache[T]) Set(reqToken, mutationToken string, value *Result[T]) {
	if value.expiredAt.IsZero() {
		value.expiredAt = def.DefaultCacheExpired()
	} else if value.expiredAt != NeverExpired && value.expiredAt.Before(time.Now()) {
		logger.Errorf("cache expired at: %v", value.expiredAt)
		return
	}

	if reqToken != "-" {
		logger.Debugf("store %s to cache: %v", reqToken, value)
		m.cacheMap.Store(reqToken, value)
	}
	if mutationToken != "-" {
		logger.Debugf("store %s to cache: %v", mutationToken, value)
		m.cacheMap.Store(mutationToken, value)
	}
}

// GetCache 获取缓存,返回 value 和 exist_state
func (m *MemoryCache[T]) GetCache(key string) (T, bool) {
	v, ok := m.cacheMap.Load(key)
	logger.Debugf("load %s from cache: %v", key, v)
	if !ok {
		return *new(T), false
	}
	if ok && v.(*Result[T]).expiredAt.Before(time.Now()) {
		m.cacheMap.Delete(key)
		return v.(*Result[T]).value, false
	}
	return v.(*Result[T]).value, ok
}

func (m *MemoryCache[T]) Delete(key string) {
	m.cacheMap.Delete(key)
}

type Result[T any] struct {
	expiredAt time.Time
	value     T
}

func NewMemoryCache[T any]() *MemoryCache[T] {
	retVal := new(MemoryCache[T])
	regClearCache(retVal)
	return retVal
}
func NewStoreResult[T any](expiredAt time.Time, value T) *Result[T] {
	return &Result[T]{
		expiredAt: expiredAt,
		value:     value,
	}
}
