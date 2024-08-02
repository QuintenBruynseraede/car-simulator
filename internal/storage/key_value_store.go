package storage

import (
	"errors"
	"fmt"
	"sync"
)

type KeyValueStoreClient struct {
	values *sync.Map
}

// NewKeyValueStoreClient returns a key-value store client initialized with an empty map
func NewKeyValueStoreClient() KeyValueStoreClient {
	return KeyValueStoreClient{
		values: new(sync.Map),
	}
}

// Read returns the value stored by `key`, and an error if the value was not found or no string
func (client *KeyValueStoreClient) Read(key string) (string, error) {
	value, ok := client.values.Load(any(key))
	if !ok {
		return "", errors.New(fmt.Sprintf("No value found for key %s", key))
	}

	strValue, ok := value.(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("Value found for key %s is not a string", key))
	}

	return strValue, nil
}

// Write inserts a `key`-`value` pair into the map
func (client *KeyValueStoreClient) Write(key string, value string) {
	client.values.Store(key, value)
}
