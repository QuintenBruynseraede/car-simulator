package controllers

import (
	"context"
	"time"

	"github.com/user/car-simulator/internal/events"
	"github.com/user/car-simulator/internal/storage"
	"go.uber.org/zap"
)

// DummyController is a controller that emits test events
type DummyController struct{}

func NewDummyController() *DummyController {
	return &DummyController{}
}

func (c *DummyController) Run(_ context.Context, kvs storage.StorageBackend, eventBus *events.EventBus, logger *zap.Logger) error {
	logger.Info("Starting Dummy controller")

	on := false
	for {
		time.Sleep(time.Second * 1)
		if on == true {
			eventBus.Publish(events.NewEvent(events.EventIndicatorLeftOff, nil))
			eventBus.Publish(events.NewEvent(events.EventIndicatorRightOff, nil))
			on = false
		} else {
			eventBus.Publish(events.NewEvent(events.EventIndicatorLeftOn, nil))
			eventBus.Publish(events.NewEvent(events.EventIndicatorRightOn, nil))
			on = true
		}
	}
	return nil
}
