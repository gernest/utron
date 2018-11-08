package cache

import (
	"errors"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

// REDIS_NIL_ERROR_RESPONSE go-redis nil response error
const REDIS_NIL_ERROR_RESPONSE = "redis: nil"

// RedisStore is the representation of the redis caching store
type RedisStore struct {
	Client redis.Client
	Prefix string
}

// Get gets a value from the store
func (rs *RedisStore) Get(key string) (interface{}, error) {
	intVal, err := rs.get(key).Int64()

	if err != nil {
		floatVal, err := rs.get(key).Float64()

		if err != nil {
			value, err := rs.get(key).Result()

			if err != nil {
				if err.Error() == REDIS_NIL_ERROR_RESPONSE {
					return "", nil
				}

				return value, err
			}

			return SimpleDecode(value)
		}

		if &floatVal == nil {
			return floatVal, errors.New("Float value is nil.")
		}

		return floatVal, nil
	}

	if &intVal == nil {
		return intVal, errors.New("Int value is nil.")
	}

	return intVal, nil
}

// GetFloat gets a float value from the store
func (rs *RedisStore) GetFloat(key string) (float64, error) {
	return rs.get(key).Float64()
}

// GetInt gets an int value from the store
func (rs *RedisStore) GetInt(key string) (int64, error) {
	return rs.get(key).Int64()
}

// Increment increments an integer counter by a given value
func (rs *RedisStore) Increment(key string, value int64) (int64, error) {
	return rs.Client.IncrBy(rs.Prefix+key, value).Result()
}

// Decrement decrements an integer counter by a given value
func (rs *RedisStore) Decrement(key string, value int64) (int64, error) {
	return rs.Client.DecrBy(rs.Prefix+key, value).Result()
}

// Put puts a value in the given store for a predetermined amount of time in mins.
func (rs *RedisStore) Put(key string, value interface{}, minutes int) error {
	time, err := time.ParseDuration(strconv.Itoa(minutes) + "m")

	if err != nil {
		return err
	}

	if IsNumeric(value) {
		return rs.Client.Set(rs.Prefix+key, value, time).Err()
	}

	val, err := Encode(value)

	if err != nil {
		return err
	}

	return rs.Client.Set(rs.GetPrefix()+key, val, time).Err()
}

// Forever puts a value in the given store until it is forgotten/evicted
func (rs *RedisStore) Forever(key string, value interface{}) error {
	if IsNumeric(value) {
		err := rs.Client.Set(rs.Prefix+key, value, 0).Err()

		if err != nil {
			return err
		}

		return rs.Client.Persist(rs.Prefix + key).Err()
	}

	val, err := Encode(value)

	if err != nil {
		return err
	}

	err = rs.Client.Set(rs.Prefix+key, val, 0).Err()

	if err != nil {
		return err
	}

	return rs.Client.Persist(rs.Prefix + key).Err()
}

// Flush flushes the store
func (rs *RedisStore) Flush() (bool, error) {
	err := rs.Client.FlushDB().Err()

	if err != nil {
		return false, err
	}

	return true, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (rs *RedisStore) Forget(key string) (bool, error) {
	err := rs.Client.Del(rs.Prefix + key).Err()

	if err != nil {
		return false, err
	}

	return true, nil
}

// GetPrefix gets the cache key prefix
func (rs *RedisStore) GetPrefix() string {
	return rs.Prefix
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (rs *RedisStore) PutMany(values map[string]interface{}, minutes int) error {
	pipe := rs.Client.TxPipeline()

	for key, value := range values {
		err := rs.Put(key, value, minutes)

		if err != nil {
			return err
		}
	}

	_, err := pipe.Exec()

	return err
}

// Many gets many values from the store
func (rs *RedisStore) Many(keys []string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	pipe := rs.Client.TxPipeline()

	for _, key := range keys {
		val, err := rs.Get(key)

		if err != nil {
			return values, err
		}

		values[key] = val
	}

	_, err := pipe.Exec()

	return values, err
}

// Connection returns the the store's client
func (rs *RedisStore) Connection() interface{} {
	return rs.Client
}

// Lpush runs the Redis lpush command
func (rs *RedisStore) Lpush(segment string, key string) {
	rs.Client.LPush(segment, key)
}

// Lrange runs the Redis lrange command
func (rs *RedisStore) Lrange(key string, start int64, stop int64) []string {
	return rs.Client.LRange(key, start, stop).Val()
}

// Tags returns the TaggedCache for the given store
func (rs *RedisStore) Tags(names ...string) TaggedStoreInterface {
	return &RedisTaggedCache{
		TaggedCache{
			Store: rs,
			Tags: TagSet{
				Store: rs,
				Names: names,
			},
		},
	}
}

// GetStruct gets the struct representation of a value from the store
func (rs *RedisStore) GetStruct(key string, entity interface{}) (interface{}, error) {
	value, err := rs.get(key).Result()

	if err != nil {
		return value, err
	}

	return Decode(value, entity)
}

func (rs *RedisStore) get(key string) *redis.StringCmd {
	return rs.Client.Get(rs.Prefix + key)
}
