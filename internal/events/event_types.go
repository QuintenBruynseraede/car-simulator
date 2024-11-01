package events

const (
	// Indicator lights
	EventIndicatorLeftOn   = "INDICATOR_LEFT_ON"
	EventIndicatorLeftOff  = "INDICATOR_LEFT_OFF"
	EventIndicatorRightOn  = "INDICATOR_RIGHT_ON"
	EventIndicatorRightOff = "INDICATOR_RIGHT_OFF"

	EventWarningEngineOn  = "EVENT_WARNING_ENGINE_ON"
	EventWarningEngineOFF = "EVENT_WARNING_ENGINE_OFF"

	EventWarningBatteryOn  = "EVENT_WARNING_BATTERY_ON"
	EventWarningBatteryOFF = "EVENT_WARNING_BATTERY_OFF"

	EventEngineStartPressed   = "EVENT_ENGINE_START_PRESSED"
	EventToggleHazardsPressed = "EVENT_TOGGLE_HAZARDS_PRESSED"
)
