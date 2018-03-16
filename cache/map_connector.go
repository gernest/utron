package cache

import (
	"errors"
)

// MapConnector is a representation of the array store connector
type MapConnector struct{}

// Connect is responsible for connecting with the caching store
func (ac *MapConnector) Connect(params map[string]interface{}) (StoreInterface, error) {
	params, err := ac.validate(params)

	if err != nil {
		return &MapStore{}, err
	}

	prefix := params["prefix"].(string)

	delete(params, "prefix")

	return &MapStore{
		Client: make(map[string]interface{}),
		Prefix: prefix,
	}, nil
}

func (ac *MapConnector) validate(params map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := params["prefix"]; !ok {
		return params, errors.New("You need to specify a caching prefix.")
	}

	return params, nil
}
