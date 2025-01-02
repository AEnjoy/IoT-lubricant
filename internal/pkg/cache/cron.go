package cache

import "github.com/AEnjoy/IoT-lubricant/pkg/utils/crontab"

func regClearCache[T any](m *MemoryCache[T]) {
	crontab.RegisterCron(m.cleanExpired, "@every 5m")
}
