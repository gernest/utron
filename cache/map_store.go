package cache

import (
	"errors"
	"strconv"
)

// MapStore is the representation of a map caching store
type MapStore struct {
	Client map[string]interface{}
	Prefix string
}

// Get gets a value from the store
func (ms *MapStore) Get(key string) (interface{}, error) {
	value := ms.Client[ms.GetPrefix()+key]

	if value == nil {
		return "", nil
	}

	if IsStringNumeric(value.(string)) {
		floatValue, err := StringToFloat64(value.(string))

		if err != nil {
			return floatValue, err
		}

		if IsFloat(floatValue) {
			return floatValue, err
		}

		return int64(floatValue), err
	}

	return SimpleDecode(value.(string))
}

// GetFloat gets a float value from the store
func (ms *MapStore) GetFloat(key string) (float64, error) {
	value := ms.Client[ms.GetPrefix()+key]

	if value == nil || !IsStringNumeric(value.(string)) {
		return 0, errors.New("Invalid numeric value")
	}

	return StringToFloat64(value.(string))
}

// GetInt gets an int value from the store
func (ms *MapStore) GetInt(key string) (int64, error) {
	value := ms.Client[ms.GetPrefix()+key]

	if value == nil || !IsStringNumeric(value.(string)) {
		return 0, errors.New("Invalid numeric value")
	}

	val, err := StringToFloat64(value.(string))

	return int64(val), err
}

// Increment increments an integer counter by a given value
func (ms *MapStore) Increment(key string, value int64) (int64, error) {
	val := ms.Client[ms.GetPrefix()+key]

	if val != nil {
		if IsStringNumeric(val.(string)) {
			floatValue, err := StringToFloat64(val.(string))

			if err != nil {
				return 0, err
			}

			result := value + int64(floatValue)

			err = ms.Put(key, result, 0)

			return result, err
		}

	}

	err := ms.Put(key, value, 0)

	return value, err
}

// Decrement decrements an integer counter by a given value
func (ms *MapStore) Decrement(key string, value int64) (int64, error) {
	return ms.Increment(key, -value)
}

// Put puts a value in the given store for a predetermined amount of time in mins.
func (ms *MapStore) Put(key string, value interface{}, minutes int) error {
	val, err := Encode(value)

	mins := strconv.Itoa(minutes)

	mins = ""

	ms.Client[ms.GetPrefix()+key+mins] = val

	return err
}

// Forever puts a value in the given store until it is forgotten/evicted
func (ms *MapStore) Forever(key string, value interface{}) error {
	return ms.Put(key, value, 0)
}

// Flush flushes the store
func (ms *MapStore) Flush() (bool, error) {
	ms.Client = make(map[string]interface{})

	return true, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (ms *MapStore) Forget(key string) (bool, error) {
	_, ok := ms.Client[ms.GetPrefix()+key]

	if ok {
		delete(ms.Client, ms.GetPrefix()+key)

		return true, nil
	}

	return false, nil
}

// GetPrefix gets the cache key prefix
func (ms *MapStore) GetPrefix() string {
	return ms.Prefix
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (ms *MapStore) PutMany(values map[string]interface{}, minutes int) error {
	for key, value := range values {
		ms.Put(key, value, minutes)
	}

	return nil
}

// Many gets many values from the store
func (ms *MapStore) Many(keys []string) (map[string]interface{}, error) {
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

// GetStruct gets the struct representation of a value from the store
func (ms *MapStore) GetStruct(key string, entity interface{}) (interface{}, error) {
	value := ms.Client[ms.GetPrefix()+key]

	return Decode(value.(string), entity)
}

// Tags returns the TaggedCache for the given store
func (ms *MapStore) Tags(names ...string) TaggedStoreInterface {
	return &TaggedCache{
		Store: ms,
		Tags: TagSet{
			Store: ms,
			Names: names,
		},
	}
}
