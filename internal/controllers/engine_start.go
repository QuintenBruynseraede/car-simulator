package controllers

import (
	"context"

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
	KeyGear     = "gear"
)

type EngineStartController struct{}

func NewEngineStartController() *EngineStartController {
	return &EngineStartController{}
}

func (c *EngineStartController) Init(kvs storage.StorageBackend, eventBus *events.EventBus, logger *zap.Logger) error {
	kvs.Write(KeyBatteryTemperature, 15.0)
	kvs.Write(KeyEngineTemperature, 15.0)
	kvs.Write(KeyIndicatorLeftStatus, "off")
	kvs.Write(KeyIndicatorRightStatus, "off")
	kvs.Write(KeyVelocity, 0.0)
	kvs.Write(KeyRPM, 0.0)
	kvs.Write(KeyEngineOn, "false")
	kvs.Write(KeyGear, "1")
	return nil
}

func (c *EngineStartController) Run(_ context.Context, kvs storage.StorageBackend, eventBus *events.EventBus, logger *zap.Logger) error {
	start_event := make(chan events.Event)
	eventBus.Subscribe(events.EventEngineStartPressed, start_event)

	for {
		select {
		case <-start_event:
			on, err := kvs.ReadString(KeyEngineOn)
			if err != nil {
				panic(err)
			}
			if on != "true" {
				logger.Info("Starting Engine")
				kvs.Write(KeyBatteryTemperature, 25.0)
				kvs.Write(KeyEngineTemperature, 25.0)
				kvs.Write(KeyVelocity, 0.0)
				kvs.Write(KeyRPM, 1000.0)
				kvs.Write(KeyEngineOn, "true")

				kvs.PauseValidation() // No validation between switching left/right off
				kvs.Write(KeyIndicatorLeftStatus, "off")
				kvs.Write(KeyIndicatorRightStatus, "off")
				kvs.StartValidation()

				logger.Info("Engine started")
			} else {
				logger.Info("Shutting down engine")
				kvs.Write(KeyRPM, 0.0)
				kvs.Write(KeyEngineOn, "false")
			}
		default:
			// Nothing
			// time.Sleep(time.Millisecond * 100)
		}
	}
}
