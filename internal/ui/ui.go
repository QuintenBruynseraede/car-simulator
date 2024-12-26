package ui

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/user/car-simulator/internal/events"
	"github.com/user/car-simulator/internal/storage"
	"go.uber.org/zap"
)

const (
	// Define offsets for UI areas
	window_size_x = 1024
	window_size_y = 768

	offset_state_x, offset_state_y       = 0, 0
	offset_controls_x, offset_controls_y = window_size_x / 2, 0
	offset_dash_x, offset_dash_y         = 0, window_size_y / 2
)

var (
	// Colors
	light_gray = color.RGBA{100, 100, 100, 255}
	gray       = color.RGBA{150, 150, 150, 255}
	black      = color.RGBA{255, 255, 255, 255}
	white      = color.RGBA{0, 0, 0, 255}
)

type UI struct {
	objects []UIObject
	logger  *zap.Logger
	state   storage.StorageBackend

	objects_under_mouse map[UIObject]bool
}

func NewUI(logger *zap.Logger, state storage.StorageBackend, event_bus *events.EventBus) *UI {
	objects := []UIObject{
		NewIndicatorLeft(),
		NewIndicatorRight(),
		NewEngineWarningLight(),
		NewBatteryWarningLight(),
		NewSpedometer(),
		NewRPMGauge(),
		NewButton("Engine On/Off", 528, 16, 100, 32, func() { event_bus.Publish(events.NewEvent(events.EventEngineStartPressed, nil)) }),
		NewButton("Toggle Hazards", 528, 58, 100, 32, func() { event_bus.Publish(events.NewEvent(events.EventToggleHazardsPressed, nil)) }),
		NewButton("Overheat engine", 528, 100, 108, 32, func() { state.Write("engine_temperature_c", 150.0) }),
		NewButton("Overheat battery", 528, 142, 114, 32, func() { state.Write("battery_temperature_c", 150.0) }),
		NewGasPedal(
			func() { event_bus.Publish(events.NewEvent(events.EventGasPedalPressed, nil)) },  // Down
			func() { event_bus.Publish(events.NewEvent(events.EventGasPedalReleased, nil)) }, // Up
		),
		NewBrakePedal(
			func() { event_bus.Publish(events.NewEvent(events.EventBrakePedalPressed, nil)) },  // Down
			func() { event_bus.Publish(events.NewEvent(events.EventBrakePedalReleased, nil)) }, // Up
		),
		NewEightSegmentDisplay(508, 416),
	}

	return &UI{objects: objects, logger: logger, state: state, objects_under_mouse: make(map[UIObject]bool)}
}

func (ui *UI) Start() error {
	ebiten.SetWindowSize(window_size_x, window_size_y)
	ebiten.SetWindowTitle("Test")
	ebiten.SetVsyncEnabled(true)

	if err := ebiten.RunGame(ui); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (ui *UI) Update() error {
	x, y := ebiten.CursorPosition()

	// Mouse hover
	for obj := range ui.objects_under_mouse {
		if !obj.IsWithinBounds(x, y) {
			obj.HandleMouseLeave(x, y)
			delete(ui.objects_under_mouse, obj)
		}
	}

	for _, obj := range ui.objects {
		if obj.IsWithinBounds(x, y) && !ui.objects_under_mouse[obj] {
			obj.HandleMouseEnter(x, y)
			ui.objects_under_mouse[obj] = true
		}
	}

	// Mouse events
	for _, obj := range ui.objects {
		if obj.IsWithinBounds(x, y) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				obj.HandleLeftClick(x, y)
			}
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				obj.HandleLeftDown(x, y)
			}
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				obj.HandleLeftUp(x, y)
			}
		}
	}

	// Updates
	for _, obj := range ui.objects {
		obj.Update(ui.state) // TODO: return copy so that they can not modify
	}

	return nil
}

func (ui *UI) Draw(screen *ebiten.Image) {
	// Draw UI area separators
	line_color := color.Gray{Y: 128}
	vector.StrokeLine(screen, offset_controls_x, offset_controls_y, offset_controls_x, offset_dash_y, 2, line_color, false)
	vector.StrokeLine(screen, offset_dash_x, offset_dash_y, window_size_x, offset_dash_y, 2, line_color, false)
	for _, obj := range ui.objects {
		obj.Draw(screen)
	}

	ebitenutil.DebugPrint(screen, stateToString(ui.state.Dump()))
}

func (ui *UI) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return window_size_x, window_size_y
}

func stateToString(state map[string]string) string {
	var buffer bytes.Buffer

	keys := make([]string, 0, len(state))
	for k := range state {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, v := range keys {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", v, state[v]))
	}

	return buffer.String()
}
