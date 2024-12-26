package internal

import (
	"context"
	"sync"

	"github.com/user/car-simulator/internal/controllers"
	"github.com/user/car-simulator/internal/events"
	"github.com/user/car-simulator/internal/logging"
	"github.com/user/car-simulator/internal/storage"
	"github.com/user/car-simulator/internal/ui"
	"go.uber.org/zap"
)

type Application struct {
	storage storage.StorageBackend
	logger  *zap.Logger
	// init_controller controllers.Controller
	controllers []controllers.Controller
	event_bus   *events.EventBus
	ui          *ui.UI
}

func NewApplication() *Application {
	logger, err := logging.NewAtLevel("info")
	storage := storage.NewKeyValueStoreClient()
	event_bus := events.NewEventBus()

	if err != nil {
		panic(err)
	}

	return &Application{
		storage: storage,
		logger:  logger,
		// init_controller:
		controllers: []controllers.Controller{
			controllers.NewEngineStartController(),
			controllers.NewIndicatorController(),
			controllers.NewPhysicsController(),
			// controllers.NewDummyController(),  // Test indicators
			// controllers.NewChaosMonkeyController(),  // Randomly changes car temperatures
		},
		event_bus: event_bus,
		ui:        ui.NewUI(logger, storage, event_bus),
	}
}

func (app *Application) Run() error {
	ctx := context.Background()

	init := sync.WaitGroup{}
	// init.Add(1)
	// go func() {
	// 	app.init_controller.Run(ctx, app.storage, app.event_bus, app.logger)
	// 	init.Done()
	// }()

	init.Wait()
	for _, controller := range app.controllers {
		go func() {
			controller.Run(ctx, app.storage, app.event_bus, app.logger)
		}()
	}

	if err := app.ui.Start(); err != nil {
		return err
	}

	return nil
}
