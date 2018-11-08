package cache

import (
	"strings"
)

// REDIS_DRIVER specifies the redis driver name
const REDIS_DRIVER = "redis"

// MEMCACHE_DRIVER specifies the memcache driver name
const MEMCACHE_DRIVER = "memcache"

// MAP_DRIVER specifies the map driver name
const MAP_DRIVER = "map"

// New new-ups an instance of StoreInterface
func New(driver string, params map[string]interface{}) (StoreInterface, error) {
	switch strings.ToLower(driver) {
	case REDIS_DRIVER:
		return connect(new(RedisConnector), params)
	case MEMCACHE_DRIVER:
		return connect(new(MemcacheConnector), params)
	case MAP_DRIVER:
		return connect(new(MapConnector), params)
	}

	return connect(new(MapConnector), params)
}

func connect(connector CacheConnectorInterface, params map[string]interface{}) (StoreInterface, error) {
	return connector.Connect(params)
}
