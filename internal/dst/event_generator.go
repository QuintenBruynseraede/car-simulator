package dst

import (
	"math/rand/v2"
	"time"

	"github.com/user/car-simulator/internal/events"
	"go.uber.org/zap"
)

var AllEvents = []events.EventType{
	events.EventEngineStartPressed,
	events.EventToggleHazardsPressed,
	events.EventGasPedalPressed,
	events.EventGasPedalReleased,
	events.EventBrakePedalPressed,
	events.EventBrakePedalReleased,
}

type EventGenerator interface {
	Run() error
}

type (
	// EventGenerator that generates no events for UI apps
	DummyEventGenerator struct{}
	// EventGenerator that generates random events on a fixed schedule
	RandomEventGenerator struct {
		logger   *zap.Logger
		eventBus *events.EventBus
		delay    time.Duration
	}
)

func NewDummyEventGenerator() EventGenerator { return DummyEventGenerator{} }
func (DummyEventGenerator) Run() error {
	return nil
}

func NewRandomEventGenerator(logger *zap.Logger, eventBus *events.EventBus) EventGenerator {
	return RandomEventGenerator{logger: logger, eventBus: eventBus, delay: time.Millisecond}
}

func (g RandomEventGenerator) Run() error {
	g.logger.Info("Starting event generator", zap.Duration("delay", g.delay))
	go func() {
		for {
			eventType := choose(AllEvents)
			g.logger.Info("Published event", zap.String("eventType", eventType))
			g.eventBus.Publish(events.NewEvent(eventType, nil))
			time.Sleep(g.delay)
		}
	}()
	return nil
}

func choose(events []events.EventType) events.EventType {
	return events[rand.Int()%len(events)]
}
