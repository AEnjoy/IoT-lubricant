package cache

import "github.com/aenjoy/iot-lubricant/pkg/utils/crontab"

func regClearCache[T any](m *MemoryCache[T]) {
	_ = crontab.RegisterCron(m.cleanExpired, "@every 5m")
}
