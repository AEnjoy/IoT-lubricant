package constant

import "time"

var DefaultCacheExpired = func() time.Time {
	return time.Now().Add(time.Hour * 6)
}
