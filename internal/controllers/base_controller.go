package controllers

import (
	"context"

	"github.com/user/car-simulator/internal/events"
	"github.com/user/car-simulator/internal/storage"
	"go.uber.org/zap"
)

type Controller interface {
	Run(ctx context.Context, kvs storage.StorageBackend, eventBus *events.EventBus, logger *zap.Logger) error
}
