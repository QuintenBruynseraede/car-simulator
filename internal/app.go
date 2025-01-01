package internal

import (
	"context"

	"github.com/user/car-simulator/internal/controllers"
	"github.com/user/car-simulator/internal/dst"
	"github.com/user/car-simulator/internal/events"
	"github.com/user/car-simulator/internal/logging"
	"github.com/user/car-simulator/internal/storage"
	"github.com/user/car-simulator/internal/ui"
	"go.uber.org/zap"
)

type Application struct {
	Storage     storage.StorageBackend
	logger      *zap.Logger
	eventBus    *events.EventBus
	ui          *ui.UI
	controllers []controllers.Controller
	// DST
	validator      dst.Validator
	simulate       bool
	validateCh     chan bool
	eventGenerator dst.EventGenerator
}

func NewUIApplication() *Application {
	logger, err := logging.NewAtLevel("info")
	if err != nil {
		panic(err)
	}

	logger.Info("Running in UI mode")
	kvs := storage.NewKeyValueStoreClient(nil)
	eventBus := events.NewEventBus()

	return &Application{
		Storage: kvs,
		logger:  logger,
		controllers: []controllers.Controller{
			controllers.NewEngineStartController(),
			controllers.NewIndicatorController(),
			controllers.NewPhysicsController(),
		},
		eventBus: eventBus,
		ui:       ui.NewUI(logger, kvs, eventBus),
		// DST
		validator:      &dst.RuntimeValidator{},
		simulate:       false,
		validateCh:     nil,
		eventGenerator: dst.NewDummyEventGenerator(),
	}
}

func NewSimulation() *Application {
	logger, err := logging.NewAtLevel("info")
	if err != nil {
		panic(err)
	}

	logger.Info("Running in headless mode")
	validateCh := make(chan bool) // Unbuffered!
	eventBus := events.NewEventBus()

	return &Application{
		Storage: storage.NewKeyValueStoreClient(validateCh),
		logger:  logger,
		controllers: []controllers.Controller{
			controllers.NewEngineStartController(),
			controllers.NewIndicatorController(),
			controllers.NewPhysicsController(),
		},
		eventBus: eventBus,
		ui:       &ui.UI{},
		// DST
		validator:      &dst.DSTValidator{Logger: logger},
		simulate:       true,
		validateCh:     validateCh,
		eventGenerator: dst.NewRandomEventGenerator(logger, eventBus),
	}
}

func (app *Application) Run() error {
	ctx := context.Background()

	// Init all controllers
	app.logger.Debug("Starting controller initialization")
	for _, controller := range app.controllers {
		controller.Init(app.Storage, app.eventBus, app.logger)
	}
	app.logger.Debug("All controllers initialized")

	// Start controllers
	for _, controller := range app.controllers {
		go func() {
			controller.Run(ctx, app.Storage, app.eventBus, app.logger)
		}()
	}

	var err error
	if app.simulate { // Run simulation
		// Move to app initialization
		app.eventGenerator.Run()
		app.Storage.StartValidation()

		for {
			msg := <-app.validateCh
			if msg {
				err = app.validator.Validate(app.Storage)
				if err != nil {
					return err
				}
			}
		}
	} else {
		err = app.ui.Start() // Run UI
		if err != nil {
			return err
		}
	}

	return nil
}
