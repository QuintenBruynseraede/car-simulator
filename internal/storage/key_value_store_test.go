package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVSCRUD(t *testing.T) {
	client := NewKeyValueStoreClient(nil)

	// Reading a key from an empty map will panic
	assert.Panics(t, func() { client.ReadString("doesntexist") })

	// Write then read returns the value
	value, key := "world", "hello"
	client.Write(key, value)
	read, err := client.ReadString(key)
	assert.Equal(t, value, read)
	assert.Nil(t, err)
}
