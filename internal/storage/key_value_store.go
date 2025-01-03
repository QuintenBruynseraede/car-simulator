package storage

import (
	"fmt"
	"sync"
)

type KeyValueStoreClient struct {
	values           *sync.Map
	validateCh       chan bool // Must be unbuffered, e.g. size 0!
	enableValidation bool
}

var (
	ErrValueNotFound  = fmt.Errorf("no value found for key")
	ErrValueNotString = fmt.Errorf("value found is not a string")
)

// NewKeyValueStoreClient returns a key-value store client initialized with an empty map
func NewKeyValueStoreClient(validateCh chan bool) *KeyValueStoreClient {
	return &KeyValueStoreClient{
		values:           new(sync.Map),
		validateCh:       validateCh,
		enableValidation: false,
	}
}

func (client *KeyValueStoreClient) StartValidation() {
	client.enableValidation = true
}

func (client *KeyValueStoreClient) PauseValidation() {
	client.enableValidation = false
}

// Read returns the value stored by `key` and panics if the value was not found or not a string
func (client *KeyValueStoreClient) ReadString(key string) (string, error) {
	value, ok := client.values.Load(any(key))
	if !ok {
		return "", ErrValueNotFound
	}

	strValue, ok := value.(string)
	if !ok {
		panic(fmt.Sprintf("Value for key %s is not a string", key))
	}

	return strValue, nil
}

// ReadFloat64 returns the value stored by `key` and panics if the value was not found or not a float
func (client *KeyValueStoreClient) ReadFloat64(key string) (float64, error) {
	value, ok := client.values.Load(any(key))
	if !ok {
		return 0, ErrValueNotFound
	}

	floatValue, ok := value.(float64)
	if !ok {
		panic(fmt.Sprintf("Value for key %s is not a float", key))
	}

	return floatValue, nil
}

// ReadInt returns the value stored by `key` and panics if the value was not found or not an int
func (client *KeyValueStoreClient) ReadInt(key string) (int, error) {
	value, ok := client.values.Load(any(key))
	if !ok {
		return 0, ErrValueNotFound
	}

	intValue, ok := value.(int)
	if !ok {
		panic(fmt.Sprintf("Value for key %s is not an int", key))
	}

	return intValue, nil
}

// Write inserts a `key`-`value` pair into the map
func (client *KeyValueStoreClient) Write(key string, value any) {
	client.values.Store(key, value)
	if client.validateCh != nil && client.enableValidation {
		client.validateCh <- true // Notify for validation and wait until it's done
	}
}

func (client *KeyValueStoreClient) Dump() map[string]string {
	dump := make(map[string]string)

	client.values.Range(func(key, value any) bool {
		switch value.(type) {
		case float64:
			dump[key.(string)] = fmt.Sprintf("%.2f", value)
		default:
			dump[key.(string)] = fmt.Sprintf("%v", value)
		}
		return true
	})

	return dump
}
