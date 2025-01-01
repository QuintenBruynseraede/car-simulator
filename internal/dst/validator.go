package dst

import (
	"errors"
	"fmt"

	"github.com/user/car-simulator/internal/storage"
	"go.uber.org/zap"
)

type Validator interface {
	Validate(state storage.StorageBackend) error
}

// Used during actual runs, doesn't validate
type RuntimeValidator struct{}

func (v *RuntimeValidator) Validate(state storage.StorageBackend) error {
	return nil
}

// Actual validator
type DSTValidator struct {
	Logger  *zap.Logger
	counter int // Number of successful validations
}

func (v *DSTValidator) Validate(state storage.StorageBackend) error {
	// velocity_kms must not be negative
	floatVal, err := state.ReadFloat64("velocity_kmh")
	if !errors.Is(err, storage.ErrValueNotFound) && floatVal < 0 {
		return DSTFailure(fmt.Sprintf("velocity_kmh is negative (%v)", floatVal), &state)
	}

	// indicator_left_status must be 'on' or 'off'
	strVal, err := state.ReadString("indicator_left_status")
	if !errors.Is(err, storage.ErrValueNotFound) && strVal != "on" && strVal != "off" {
		return DSTFailure(fmt.Sprintf("indicator_left_status is not on or off (%v)", strVal), &state)
	}

	// engine_rpm must be between 0 and 8000
	floatVal, err = state.ReadFloat64("engine_rpm")
	if !errors.Is(err, storage.ErrValueNotFound) && (floatVal < 0 || floatVal > 8000) {
		return DSTFailure(fmt.Sprintf("engine_rpm is not between 0 and 8000 (%v)", floatVal), &state)
	}

	// Both indicators must have the same status
	left, err := state.ReadString("indicator_left_status")
	right, err := state.ReadString("indicator_right_status")
	if !errors.Is(err, storage.ErrValueNotFound) && left != right {
		return DSTFailure(fmt.Sprintf("indicators have different statuses (%v, %v)", left, right), &state)
	}

	v.counter++
	if v.counter%1e6 == 0 {
		v.Logger.Info("Validation successful", zap.Int("counter", v.counter))
	}
	return nil
}
