package events

type EventType = string

const (
	EventEngineStartPressed   = "EVENT_ENGINE_START_PRESSED"
	EventToggleHazardsPressed = "EVENT_TOGGLE_HAZARDS_PRESSED"

	EventGasPedalPressed    = "EVENT_GAS_PEDAL_PRESSED"
	EventGasPedalReleased   = "EVENT_GAS_PEDAL_RELEASED"
	EventBrakePedalPressed  = "EVENT_BRAKE_PEDAL_PRESSED"
	EventBrakePedalReleased = "EVENT_BRAKE_PEDAL_RELEASED"
)
