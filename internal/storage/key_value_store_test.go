package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVSCRUD(t *testing.T) {
	client := NewKeyValueStoreClient()

	// Reading a key from an empty map returns an error
	value, error := client.Read("doesntexist")
	assert.Error(t, error)
	assert.Equal(t, "", value)

	// Write then read returns the value
	value = "world"
	key := "hello"
	client.Write(key, value)
	read, error := client.Read(key)
	assert.NoError(t, error)
	assert.Equal(t, value, read)
}
