package controllers

import (
	"context"
	"time"

	"github.com/user/car-simulator/internal/events"
	"github.com/user/car-simulator/internal/storage"
	"go.uber.org/zap"
)

// IndicatorController is responsible for controlling the car's indicators
type IndicatorController struct {
	hazardsOn bool
}

func NewIndicatorController() *IndicatorController {
	return &IndicatorController{
		hazardsOn: false,
	}
}

func (c *IndicatorController) Run(_ context.Context, kvs storage.StorageBackend, eventBus *events.EventBus, logger *zap.Logger) error {
	logger.Info("Starting Indicator controller")

	toggle := make(chan events.Event)
	eventBus.Subscribe(events.EventToggleHazardsPressed, toggle)

	for {
		select {
		case <-toggle:
			c.hazardsOn = !c.hazardsOn
		default:
			if c.hazardsOn {
				toggleIndicators(kvs)
			} else {
				disableIndicators(kvs)
			}
			time.Sleep(time.Millisecond * 800)
		}
	}
}

// Flip the state of both indicators
func toggleIndicators(kvs storage.StorageBackend) {
	stateRight, err := kvs.ReadString(KeyIndicatorRightStatus)
	if err != nil {
		panic(err)
	}

	if stateRight == "on" {
		kvs.Write(KeyIndicatorRightStatus, "off")
	} else {
		kvs.Write(KeyIndicatorRightStatus, "on")
	}

	stateLeft, err := kvs.ReadString(KeyIndicatorLeftStatus)
	if err != nil {
		panic(err)
	}

	if stateLeft == "on" {
		kvs.Write(KeyIndicatorLeftStatus, "off")
	} else {
		kvs.Write(KeyIndicatorLeftStatus, "on")
	}
}

func disableIndicators(kvs storage.StorageBackend) {
	kvs.Write(KeyIndicatorRightStatus, "off")
	kvs.Write(KeyIndicatorLeftStatus, "off")
}
