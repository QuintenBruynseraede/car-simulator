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
	event_bus   *events.EventBus
	ui          *ui.UI
	controllers []controllers.Controller
	validator   dst.Validator
	simulate    bool
	validateCh  chan bool
}

func NewApplication(simulate bool) *Application {
	logger, err := logging.NewAtLevel("info")
	event_bus := events.NewEventBus()

	if err != nil {
		panic(err)
	}

	var userInterface *ui.UI
	var validator dst.Validator
	var kvs storage.StorageBackend
	var validateCh chan bool

	if simulate {
		logger.Info("Running in headless mode")
		userInterface = &ui.UI{}
		validator = &dst.DSTValidator{Logger: logger}
		validateCh = make(chan bool) // Unbuffered, i.e. size 0!
		kvs = storage.NewKeyValueStoreClient(validateCh)
	} else {
		userInterface = ui.NewUI(logger, kvs, event_bus)
		validator = &dst.RuntimeValidator{}
		kvs = storage.NewKeyValueStoreClient(nil) // No validation channel
	}

	return &Application{
		Storage: kvs,
		logger:  logger,
		// init_controller:
		controllers: []controllers.Controller{
			controllers.NewEngineStartController(),
			controllers.NewIndicatorController(),
			controllers.NewPhysicsController(),
		},
		event_bus:  event_bus,
		ui:         userInterface,
		validator:  validator,
		simulate:   simulate,
		validateCh: validateCh,
	}
}

func (app *Application) Run() error {
	ctx := context.Background()

	for _, controller := range app.controllers {
		go func() {
			controller.Run(ctx, app.Storage, app.event_bus, app.logger)
		}()
	}

	var err error
	if app.simulate { // Run simulation
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
