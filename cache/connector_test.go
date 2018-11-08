package cache

import (
	"testing"
)

func TestMemcacheConnector(t *testing.T) {
	memcacheConnector := new(MemcacheConnector)

	memcacheStore, err := memcacheConnector.Connect(memcacheStore())

	if err != nil {
		panic(err)
	}

	_, ok := memcacheStore.(StoreInterface)

	if !ok {
		t.Error("Expected StoreInterface got", memcacheStore)
	}
}

func TestRedisConnector(t *testing.T) {
	redisConnector := new(RedisConnector)

	redisStore, err := redisConnector.Connect(redisStore())

	if err != nil {
		panic(err)
	}

	_, ok := redisStore.(StoreInterface)

	if !ok {
		t.Error("Expected StoreInterface got", redisStore)
	}
}

func TestArrayConnector(t *testing.T) {
	mapConnector := new(MapConnector)

	mapStore, err := mapConnector.Connect(mapStore())

	if err != nil {
		panic(err)
	}

	_, ok := mapStore.(StoreInterface)

	if !ok {
		t.Error("Expected StoreInterface got", mapStore)
	}
}

func redisStore() map[string]interface{} {
	params := make(map[string]interface{})

	params["address"] = "localhost:6379"
	params["password"] = ""
	params["database"] = 0
	params["prefix"] = "golavel:"

	return params
}

func memcacheStore() map[string]interface{} {
	params := make(map[string]interface{})

	params["server 1"] = "127.0.0.1:11211"
	params["prefix"] = "golavel:"

	return params
}

func mapStore() map[string]interface{} {
	params := make(map[string]interface{})

	params["prefix"] = "golavel:"

	return params
}
