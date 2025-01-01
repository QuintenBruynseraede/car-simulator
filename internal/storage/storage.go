package storage

type StorageBackend interface {
	// ReadString returns a string value from the state. Returns an error
	// if the value is not found. Panics if the value cannot be cast to a string.
	ReadString(key string) (string, error)
	// ReadFloat64 returns a float64 value from the state. Returns an error
	// if the value is not found. Panics if the value cannot be cast to a float64.
	ReadFloat64(key string) (float64, error)
	// ReadInt returns an integer value from the state. Returns an error
	// if the value is not found. Panics if the value cannot be cast to an integer.
	ReadInt(key string) (int, error)
	// Write stores a value in the state
	Write(key string, value any)
	// Dump returns a map which contains all the values in the state.
	// For non-string values stored, a string representation is generated
	Dump() map[string]string
	// StartValidation signals that the KVS can start sending events to the
	// validation channel to trigger state validation.
	StartValidation()
	// PauseValidation signals that validation must be temporary disabled to
	// allow for multiple writes without intermediate validation. Validation
	// must be explicitly enabled again using StartValidation().
	PauseValidation()
}
