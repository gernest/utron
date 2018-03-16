package cache

// CacheConnectorInterface represents the connector methods to be implemented
type CacheConnectorInterface interface {
	Connect(params map[string]interface{}) (StoreInterface, error)

	validate(params map[string]interface{}) (map[string]interface{}, error)
}

// CacheInterface represents the caching methods to be implemented
type CacheInterface interface {
	Get(key string) (interface{}, error)

	Put(key string, value interface{}, minutes int) error

	Increment(key string, value int64) (int64, error)

	Decrement(key string, value int64) (int64, error)

	Forget(key string) (bool, error)

	Forever(key string, value interface{}) error

	Flush() (bool, error)

	GetInt(key string) (int64, error)

	GetFloat(key string) (float64, error)

	GetPrefix() string

	Many(keys []string) (map[string]interface{}, error)

	PutMany(values map[string]interface{}, minutes int) error

	GetStruct(key string, entity interface{}) (interface{}, error)
}

// TagsInterface represents the tagging methods to be implemented
type TagsInterface interface {
	Tags(names ...string) TaggedStoreInterface
}

// StoreInterface represents the methods a caching store needs to implement
type StoreInterface interface {
	CacheInterface

	TagsInterface
}

// TaggedStoreInterface represents the methods a tagged-caching store needs to implement
type TaggedStoreInterface interface {
	CacheInterface

	TagFlush() error
}
