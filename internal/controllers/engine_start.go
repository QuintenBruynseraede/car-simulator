package controllers

import (
	"context"
	"time"

	"github.com/user/car-simulator/internal/events"
	"github.com/user/car-simulator/internal/storage"
	"go.uber.org/zap"
)

const (
	KeyEngineOn           = "engine_on"
	KeyEngineTemperature  = "engine_temperature_c"
	KeyBatteryTemperature = "battery_temperature_c"

	KeyIndicatorLeftStatus  = "indicator_left_status"
	KeyIndicatorRightStatus = "indicator_right_status"

	KeyVelocity = "velocity_kmh"
	KeyRPM      = "engine_rpm"
)

// ChaosMonkeyController is a controller that randomly modifies the car's parameters
type EngineStartController struct{}

func NewEngineStartController() *EngineStartController {
	return &EngineStartController{}
}

func (c *EngineStartController) Run(_ context.Context, kvs storage.StorageBackend, eventBus *events.EventBus, logger *zap.Logger) error {
	start_event := make(chan events.Event)
	eventBus.Subscribe(events.EventEngineStartPressed, start_event)

	kvs.Write(KeyBatteryTemperature, 15.0)
	kvs.Write(KeyEngineTemperature, 15.0)
	kvs.Write(KeyIndicatorLeftStatus, "off")
	kvs.Write(KeyIndicatorRightStatus, "off")
	kvs.Write(KeyVelocity, 0.0)
	kvs.Write(KeyRPM, 0.0)
	kvs.Write(KeyEngineOn, "false")

	for {
		select {
		case <-start_event:
			if kvs.ReadString(KeyEngineOn) != "true" {
				logger.Info("Starting Engine")
				kvs.Write(KeyBatteryTemperature, 25.0)
				kvs.Write(KeyEngineTemperature, 25.0)
				kvs.Write(KeyIndicatorLeftStatus, "off")
				kvs.Write(KeyIndicatorRightStatus, "off")
				kvs.Write(KeyVelocity, 0.0)
				kvs.Write(KeyRPM, 1000.0)
				kvs.Write(KeyEngineOn, "true")

				logger.Info("Engine started")
			} else {
				logger.Info("Shutting down engine")
				kvs.Write(KeyRPM, 0.0)
				kvs.Write(KeyEngineOn, "false")
			}
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
