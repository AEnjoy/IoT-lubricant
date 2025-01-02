package cache

import (
	"sync"
	"time"

	def "github.com/AEnjoy/IoT-lubricant/pkg/default"
)

type MemoryCache[T any] struct {
	sync.Mutex
	cacheMap map[string]*Result[T]
}

func (m *MemoryCache[T]) cleanExpired() {
	m.Lock()
	defer m.Unlock()

	for k, v := range m.cacheMap {
		if v.expiredAt.Before(time.Now()) {
			delete(m.cacheMap, k)
		}
	}
}
func (m *MemoryCache[T]) Set(reqToken, mutationToken string, value *Result[T]) {
	m.Lock()
	defer m.Unlock()

	if value.expiredAt.IsZero() {
		value.expiredAt = def.DefaultCacheExpired()
	} else if value.expiredAt.Before(time.Now()) {
		return
	}

	if reqToken != "-" {
		m.cacheMap[reqToken] = value
	}
	m.cacheMap[mutationToken] = value
}

// GetCache 获取缓存,返回 value 和 exist_state
func (m *MemoryCache[T]) GetCache(key string) (T, bool) {
	m.Lock()
	defer m.Unlock()
	v, ok := m.cacheMap[key]
	if ok && v.expiredAt.Before(time.Now()) {
		delete(m.cacheMap, key)
		return v.value, false
	}
	return v.value, ok
}
func (m *MemoryCache[T]) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.cacheMap, key)
}

type Result[T any] struct {
	expiredAt time.Time
	value     T
}

func NewMemoryCache[T any]() *MemoryCache[T] {
	retVal := &MemoryCache[T]{
		cacheMap: make(map[string]*Result[T], 50),
	}
	regClearCache(retVal)
	return retVal
}
