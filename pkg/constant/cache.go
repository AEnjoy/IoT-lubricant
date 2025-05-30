package constant

import "time"

var DefaultCacheExpired = func() time.Time {
	return time.Now().Add(time.Hour * 6)
}

const LatestDataCacheKey = "latest_data_cache_key-%s-%s" // project_id,agent_id
