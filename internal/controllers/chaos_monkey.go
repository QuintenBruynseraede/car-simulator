package controllers

import (
	"context"
	"math/rand"
	"time"

	"github.com/user/car-simulator/internal/events"
	"github.com/user/car-simulator/internal/storage"
	"go.uber.org/zap"
)

// ChaosMonkeyController is a controller that randomly modifies the car's parameters
type ChaosMonkeyController struct{}

func NewChaosMonkeyController() *ChaosMonkeyController {
	return &ChaosMonkeyController{}
}

func (c *ChaosMonkeyController) Run(_ context.Context, kvs storage.StorageBackend, eventBus *events.EventBus, logger *zap.Logger) error {
	logger.Info("Starting ChaosMonkey controller")
	rand.Seed(1234)

	time.Sleep(time.Second * 3 * time.Duration(rand.Intn(30)))
	kvs.Write(KeyBatteryTemperature, "250")

	time.Sleep(time.Second * 3 * time.Duration(rand.Intn(30)))
	kvs.Write(KeyEngineTemperature, "250")

	return nil
}
