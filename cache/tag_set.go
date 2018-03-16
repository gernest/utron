package cache

import (
	"github.com/segmentio/ksuid"
	"strings"
)

// TagSet is the representation of a tag set for the cahing stores
type TagSet struct {
	Store StoreInterface
	Names []string
}

// GetNamespace gets the current TagSet namespace
func (ts *TagSet) GetNamespace() (string, error) {
	tagsIds, err := ts.tagIds()

	if err != nil {
		return "", err
	}

	return strings.Join(tagsIds, "|"), err
}

// Reset resets the tag set
func (ts *TagSet) Reset() error {
	for i, name := range ts.Names {
		id, err := ts.resetTag(name)

		if err != nil {
			return err
		}

		ts.Names[i] = id
	}

	return nil
}

func (ts *TagSet) tagId(name string) (string, error) {
	value, err := ts.Store.Get(ts.tagKey(name))

	if err != nil {
		return value.(string), err
	}

	if value == "" {
		return ts.resetTag(name)
	}

	return value.(string), nil
}

func (ts *TagSet) tagKey(name string) string {
	return "tag:" + name + ":key"
}

func (ts *TagSet) tagIds() ([]string, error) {
	tagIds := make([]string, len(ts.Names))

	for i, name := range ts.Names {
		val, err := ts.tagId(name)

		if err != nil {
			return tagIds, err
		}

		tagIds[i] = val
	}

	return tagIds, nil
}

func (ts *TagSet) resetTag(name string) (string, error) {
	id := ksuid.New().String()

	err := ts.Store.Forever(ts.tagKey(name), id)

	return id, err
}
