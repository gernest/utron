package cache

import (
	"crypto/sha1"
	"encoding/hex"
)

// TaggedCache is the representation of a tagged caching store
type TaggedCache struct {
	Store StoreInterface
	Tags  TagSet
}

// Get gets a value from the store
func (tc *TaggedCache) Get(key string) (interface{}, error) {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return tagKey, err
	}

	return tc.Store.Get(tagKey)
}

// Put puts a value in the given store for a predetermined amount of time in mins.
func (tc *TaggedCache) Put(key string, value interface{}, minutes int) error {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return err
	}

	return tc.Store.Put(tagKey, value, minutes)
}

// Increment increments an integer counter by a given value
func (tc *TaggedCache) Increment(key string, value int64) (int64, error) {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return 0, err
	}

	return tc.Store.Increment(tagKey, value)
}

// Decrement decrements an integer counter by a given value
func (tc *TaggedCache) Decrement(key string, value int64) (int64, error) {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return 0, err
	}

	return tc.Store.Decrement(tagKey, value)
}

// Forget forgets/evicts a given key-value pair from the store
func (tc *TaggedCache) Forget(key string) (bool, error) {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return false, err
	}

	return tc.Store.Forget(tagKey)
}

// Forever puts a value in the given store until it is forgotten/evicted
func (tc *TaggedCache) Forever(key string, value interface{}) error {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return err
	}

	return tc.Store.Forever(tagKey, value)
}

// Flush flushes the store
func (tc *TaggedCache) Flush() (bool, error) {
	return tc.Store.Flush()
}

// Many gets many values from the store
func (tc *TaggedCache) Many(keys []string) (map[string]interface{}, error) {
	taggedKeys := make([]string, len(keys))
	values := make(map[string]interface{})

	for i, key := range keys {
		tagKey, err := tc.taggedItemKey(key)

		if err != nil {
			return values, err
		}

		taggedKeys[i] = tagKey
	}

	results, err := tc.Store.Many(taggedKeys)

	if err != nil {
		return results, err
	}

	for i, result := range results {
		values[GetTaggedManyKey(tc.GetPrefix(), i)] = result
	}

	return values, nil
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (tc *TaggedCache) PutMany(values map[string]interface{}, minutes int) error {
	taggedMap := make(map[string]interface{})

	for key, value := range values {
		tagKey, err := tc.taggedItemKey(key)

		if err != nil {
			return err
		}

		taggedMap[tagKey] = value
	}

	return tc.Store.PutMany(taggedMap, minutes)
}

// GetPrefix gets the cache key prefix
func (tc *TaggedCache) GetPrefix() string {
	return tc.Store.GetPrefix()
}

// GetInt gets an int value from the store
func (tc *TaggedCache) GetInt(key string) (int64, error) {
	return tc.Store.GetInt(key)
}

// GetFloat gets a float value from the store
func (tc *TaggedCache) GetFloat(key string) (float64, error) {
	return tc.Store.GetFloat(key)
}

// GetStruct gets the struct representation of a value from the store
func (tc *TaggedCache) GetStruct(key string, entity interface{}) (interface{}, error) {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return tagKey, err
	}

	return tc.Store.GetStruct(tagKey, entity)
}

// TagFlush flushes the tags of the TaggedCache
func (tc *TaggedCache) TagFlush() error {
	return tc.Tags.Reset()
}

// GetTags returns the TaggedCache Tags
func (tc *TaggedCache) GetTags() TagSet {
	return tc.Tags
}

func (tc *TaggedCache) taggedItemKey(key string) (string, error) {
	h := sha1.New()

	namespace, err := tc.Tags.GetNamespace()

	if err != nil {
		return namespace, err
	}

	h.Write(([]byte(namespace)))

	return tc.GetPrefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key, nil
}
