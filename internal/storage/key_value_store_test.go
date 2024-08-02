package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVSCRUD(t *testing.T) {
	client := NewKeyValueStoreClient()

	// Reading a key from an empty map returns an error
	value, err := client.Read("doesntexist")
	assert.Error(t, err)
	assert.Equal(t, "", value)

	// Write then read returns the value
	value = "world"
	key := "hello"
	client.Write(key, value)
	read, err := client.Read(key)
	assert.NoError(t, err)
	assert.Equal(t, value, read)
}
