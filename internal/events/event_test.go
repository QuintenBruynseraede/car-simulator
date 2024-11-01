package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventBus(t *testing.T) {
	bus := NewEventBus()
	event := NewEvent("test-event", "hello")

	// Subscribe to the event
	subscriber := make(chan Event)
	bus.Subscribe(event.Type, subscriber)

	// Publish the event
	go func() { bus.Publish(event) }()

	// Check that the subscriber received the event
	receivedEvent := <-subscriber
	assert.Equal(t, event, receivedEvent)
}
