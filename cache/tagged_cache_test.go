package cache

import "testing"

func TestPutGetWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		expected := "value"

		tags := tag()

		cache.Tags(tags).Put("key", "value", 10)

		got, err := cache.Tags(tags).Get("key")

		if got != expected || err != nil {
			t.Error("Expected value, got ", got)
		}

		cache.Tags(tags).Forget("key")
	}
}

func TestPutGetIntWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		tags := tag()

		cache.Tags(tags).Put("key", 100, 1)

		got, err := cache.Tags(tags).Get("key")

		if got != int64(100) || err != nil {
			t.Error("Expected 100, got ", got)
		}

		cache.Tags(tags).Forget("key")
	}
}

func TestPutGetFloatWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		var expected float64

		expected = 9.99

		tags := tag()

		cache.Tags(tags).Put("key", expected, 1)

		got, err := cache.Tags(tags).Get("key")

		if got != expected || err != nil {
			t.Error("Expected 9.99, got ", got)
		}

		cache.Tags(tags).Forget("key")
	}
}

func TestIncrementWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		tags := tag()

		cache.Tags(tags).Increment("increment_key", 1)
		cache.Tags(tags).Increment("increment_key", 1)
		got, err := cache.Tags(tags).Get("increment_key")

		var expected int64 = 2

		if got != expected || err != nil {
			t.Error("Expected 2, got ", got)
		}

		cache.Tags(tags).Forget("increment_key")
	}
}

func TestDecrementWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		tags := tag()

		cache.Tags(tags).Increment("decrement_key", 2)
		cache.Tags(tags).Decrement("decrement_key", 1)

		var expected int64 = 1

		got, err := cache.Tags(tags).Get("decrement_key")

		if got != expected || err != nil {
			t.Error("Expected "+string(expected)+", got ", got)
		}

		cache.Tags(tags).Forget("decrement_key")
	}
}

func TestForeverWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		expected := "value"

		tags := tag()

		cache.Tags(tags).Forever("key", expected)

		got, err := cache.Tags(tags).Get("key")

		if got != expected || err != nil {
			t.Error("Expected "+expected+", got ", got)
		}

		cache.Tags(tags).Forget("key")
	}
}

func TestPutGetManyWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		tags := tag()

		keys := make(map[string]interface{})

		keys["key_1"] = "value"
		keys["key_2"] = int64(100)
		keys["key_3"] = float64(9.99)

		cache.Tags(tags).PutMany(keys, 10)

		resultKeys := make([]string, 3)

		resultKeys[0] = "key_1"
		resultKeys[1] = "key_2"
		resultKeys[2] = "key_3"

		results, err := cache.Tags(tags).Many(resultKeys)

		if err != nil {
			panic(err)
		}

		for i := range results {
			if results[i] != keys[i] {
				t.Error(i, results[i])
			}
		}

		cache.Tags(tags).Flush()
	}
}

func TestPutGetStructWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		tags := make([]string, 3)

		tags[0] = "tag1"
		tags[1] = "tag2"
		tags[2] = "tag3"

		var example Example

		example.Name = "Alejandro"
		example.Description = "Whatever"

		cache.Tags(tags...).Put("key", example, 10)

		var newExample Example

		cache.Tags(tags...).GetStruct("key", &newExample)

		if newExample != example {
			t.Error("The structs are not the same", newExample)
		}

		cache.Forget("key")
	}
}

func TestTagSet(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		tagSet := cache.Tags("Alejandro").GetTags()

		namespace, err := tagSet.GetNamespace()

		if err != nil {
			panic(err)
		}

		if len([]rune(namespace)) != 20 {
			t.Error("The namespace is not 20 chars long.", namespace)
		}

		got := tagSet.Reset()

		if got != nil {
			t.Error("Reset did not return nil.", got)
		}
	}
}

func tag() string {
	return "tag"
}
