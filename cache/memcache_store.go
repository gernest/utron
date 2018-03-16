package cache

import (
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
)

// MEMCACHE_NIL_ERROR_RESPONSE is the gomemcache nil response error
const MEMCACHE_NIL_ERROR_RESPONSE = "memcache: cache miss"

// MemcacheStore is the representation of the memcache caching store
type MemcacheStore struct {
	Client memcache.Client
	Prefix string
}

// Put puts a value in the given store for a predetermined amount of time in mins.
func (ms *MemcacheStore) Put(key string, value interface{}, minutes int) error {
	item, err := ms.item(key, value, minutes)

	if err != nil {
		return err
	}

	return ms.Client.Set(item)
}

// Forever puts a value in the given store until it is forgotten/evicted
func (ms *MemcacheStore) Forever(key string, value interface{}) error {
	return ms.Put(key, value, 0)
}

// Get gets a value from the store
func (ms *MemcacheStore) Get(key string) (interface{}, error) {
	value, err := ms.get(key)

	if err != nil {
		return value, err
	}

	return ms.processValue(value)
}

// GetFloat gets a float value from the store
func (ms *MemcacheStore) GetFloat(key string) (float64, error) {
	value, err := ms.get(key)

	if err != nil {
		return 0.0, err
	}

	if !IsStringNumeric(value) {
		return 0.0, errors.New("Invalid numeric value")
	}

	return StringToFloat64(value)
}

// GetInt gets an int value from the store
func (ms *MemcacheStore) GetInt(key string) (int64, error) {
	value, err := ms.get(key)

	if err != nil {
		return 0, err
	}

	if !IsStringNumeric(value) {
		return 0, errors.New("Invalid numeric value")
	}

	val, err := StringToFloat64(value)

	return int64(val), err
}

// Increment increments an integer counter by a given value
func (ms *MemcacheStore) Increment(key string, value int64) (int64, error) {
	newValue, err := ms.Client.Increment(ms.GetPrefix()+key, uint64(value))

	if err != nil {
		if err.Error() != "memcache: cache miss" {
			return value, err
		}

		ms.Put(key, value, 0)

		return value, nil
	}

	return int64(newValue), nil
}

// Decrement decrements an integer counter by a given value
func (ms *MemcacheStore) Decrement(key string, value int64) (int64, error) {
	newValue, err := ms.Client.Decrement(ms.GetPrefix()+key, uint64(value))

	if err != nil {
		if err.Error() != "memcache: cache miss" {
			return value, err
		}

		ms.Put(key, 0, 0)

		return int64(0), nil
	}

	return int64(newValue), nil
}

// GetPrefix gets the cache key prefix
func (ms *MemcacheStore) GetPrefix() string {
	return ms.Prefix
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (ms *MemcacheStore) PutMany(values map[string]interface{}, minutes int) error {
	for key, value := range values {
		err := ms.Put(key, value, minutes)

		if err != nil {
			return err
		}
	}

	return nil
}

// Many gets many values from the store
func (ms *MemcacheStore) Many(keys []string) (map[string]interface{}, error) {
	items := make(map[string]interface{})

	for _, key := range keys {
		val, err := ms.Get(key)

		if err != nil {
			return items, err
		}

		items[key] = val
	}

	return items, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (ms *MemcacheStore) Forget(key string) (bool, error) {
	err := ms.Client.Delete(ms.GetPrefix() + key)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Flush flushes the store
func (ms *MemcacheStore) Flush() (bool, error) {
	err := ms.Client.DeleteAll()

	if err != nil {
		return false, err
	}

	return true, nil
}

// GetStruct gets the struct representation of a value from the store
func (ms *MemcacheStore) GetStruct(key string, entity interface{}) (interface{}, error) {
	value, err := ms.get(key)

	if err != nil {
		return value, err
	}

	return Decode(value, entity)
}

// Tags returns the TaggedCache for the given store
func (ms *MemcacheStore) Tags(names ...string) TaggedStoreInterface {
	return &TaggedCache{
		Store: ms,
		Tags: TagSet{
			Store: ms,
			Names: names,
		},
	}
}

func (ms *MemcacheStore) get(key string) (string, error) {
	item, err := ms.Client.Get(ms.GetPrefix() + key)

	if err != nil {
		if err.Error() == MEMCACHE_NIL_ERROR_RESPONSE {
			return "", nil
		}

		return "", err
	}

	return ms.getItemValue(item.Value), nil
}

func (ms *MemcacheStore) getItemValue(itemValue []byte) string {
	value, err := SimpleDecode(string(itemValue))

	if err != nil {
		return string(itemValue)
	}

	return value
}

func (ms *MemcacheStore) processValue(value string) (interface{}, error) {
	if IsStringNumeric(value) {
		floatValue, err := StringToFloat64(value)

		if err != nil {
			return floatValue, err
		}

		if IsFloat(floatValue) {
			return floatValue, err
		}

		return int64(floatValue), err
	}

	return value, nil
}

func (ms *MemcacheStore) item(key string, value interface{}, minutes int) (*memcache.Item, error) {
	val, err := Encode(value)

	return &memcache.Item{
		Key:        ms.GetPrefix() + key,
		Value:      []byte(val),
		Expiration: int32(minutes),
	}, err
}
