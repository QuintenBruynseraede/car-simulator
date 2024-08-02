package storage

import (
	"fmt"
	"sync"
)

type KeyValueStoreClient struct {
	values *sync.Map
}

var (
	ErrValueNotFound  = fmt.Errorf("No value found for key")
	ErrValueNotString = fmt.Errorf("Value found is not a string")
)

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
		return "", ErrValueNotFound
	}

	strValue, ok := value.(string)
	if !ok {
		return "", ErrValueNotString
	}

	return strValue, nil
}

// Write inserts a `key`-`value` pair into the map
func (client *KeyValueStoreClient) Write(key string, value string) {
	client.values.Store(key, value)
}
