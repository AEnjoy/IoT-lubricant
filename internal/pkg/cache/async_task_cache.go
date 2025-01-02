package cache

import (
	"sync"
	"time"

	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
)

type memoryCache struct {
	sync.Mutex
	cacheMap map[string]*Result
}
type Result struct {
	expiredAt time.Time
	value     *core.QueryTaskResultResponse
}

var mc = memoryCache{
	cacheMap: make(map[string]*Result, 50),
}

func NewResult(expiredAt time.Time, value *core.QueryTaskResultResponse) *Result {
	return &Result{
		expiredAt: expiredAt,
		value:     value,
	}
}
func SetCache(reqToken, mutationToken string, value *Result) {
	mc.Lock()
	defer mc.Unlock()
	if reqToken != "-" {
		mc.cacheMap[reqToken] = value
	}
	mc.cacheMap[mutationToken] = value
}
func GetCache(key string) (*Result, bool) {
	mc.Lock()
	defer mc.Unlock()
	v, ok := mc.cacheMap[key]
	return v, ok
}
func DelCache(key string) {
	mc.Lock()
	defer mc.Unlock()
	delete(mc.cacheMap, key)
}
