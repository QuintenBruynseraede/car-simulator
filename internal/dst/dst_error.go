package dst

import (
	"fmt"

	"github.com/user/car-simulator/internal/storage"
)

type DSTAssertionError struct {
	State   *storage.StorageBackend
	Message string
}

func (r *DSTAssertionError) Error() string {
	return fmt.Sprintf("dst assertion failed: %v: ", r.Message)
}

func DSTFailure(message string, state *storage.StorageBackend) *DSTAssertionError {
	return &DSTAssertionError{
		Message: message,
		State:   state,
	}
}
