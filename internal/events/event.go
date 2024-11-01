package events

import "time"

type Event struct {
	Type      string
	Timestamp time.Time
	Data      any
}

func NewEvent(eventType string, data any) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
}

type EventBus struct {
	subscribers map[string][]chan<- Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan<- Event),
	}
}

func (bus *EventBus) Subscribe(eventType string, subscriber chan<- Event) {
	bus.subscribers[eventType] = append(bus.subscribers[eventType], subscriber)
}

func (bus *EventBus) Publish(event Event) {
	for _, subscriber := range bus.subscribers[event.Type] {
		subscriber <- event
	}
}
