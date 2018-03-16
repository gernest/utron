package cache

import (
	"crypto/sha1"
	"encoding/hex"
	"reflect"
	"strings"
)

// RedisTaggedCache is the representation of the redis tagged cache store
type RedisTaggedCache struct {
	TaggedCache
}

// Forever puts a value in the given store until it is forgotten/evicted
func (rtc *RedisTaggedCache) Forever(key string, value interface{}) error {
	namespace, err := rtc.Tags.GetNamespace()

	if err != nil {
		return err
	}

	rtc.pushForever(namespace, key)

	h := sha1.New()

	h.Write(([]byte(namespace)))

	return rtc.Store.Forever(rtc.GetPrefix()+hex.EncodeToString(h.Sum(nil))+":"+key, value)
}

// TagFlush flushes the tags of the TaggedCache
func (rtc *RedisTaggedCache) TagFlush() error {
	return rtc.deleteForeverKeys()
}

func (rtc *RedisTaggedCache) pushForever(namespace string, key string) {
	h := sha1.New()

	h.Write(([]byte(namespace)))

	fullKey := rtc.GetPrefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key

	segments := strings.Split(namespace, "|")

	for _, segment := range segments {

		inputs := []reflect.Value{
			reflect.ValueOf(rtc.foreverKey(segment)),
			reflect.ValueOf(fullKey),
		}

		reflect.ValueOf(rtc.Store).MethodByName("Lpush").Call(inputs)
	}
}

func (rtc *RedisTaggedCache) deleteForeverKeys() error {
	namespace, err := rtc.Tags.GetNamespace()

	if err != nil {
		return err
	}

	segments := strings.Split(namespace, "|")

	for _, segment := range segments {
		key := rtc.foreverKey(segment)

		err = rtc.deleteForeverValues(key)

		if err != nil {
			return err
		}

		_, err = rtc.Store.Forget(segment)

		if err != nil {
			return err
		}
	}

	return nil
}

func (rtc *RedisTaggedCache) deleteForeverValues(key string) error {
	inputs := []reflect.Value{
		reflect.ValueOf(key),
		reflect.ValueOf(int64(0)),
		reflect.ValueOf(int64(-1)),
	}

	keys := reflect.ValueOf(rtc.Store).MethodByName("Lrange").Call(inputs)

	if len(keys) > 0 {
		for _, key := range keys {
			if key.Len() > 0 {
				for i := 0; i < key.Len(); i++ {
					_, err := rtc.Store.Forget(key.Index(i).String())

					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (rtc *RedisTaggedCache) foreverKey(segment string) string {
	return rtc.GetPrefix() + segment + ":forever"
}
