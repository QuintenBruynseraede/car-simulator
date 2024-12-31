package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/user/car-simulator/internal/storage"
)

const (
	path_engine_on   = "assets/warning_engine_on.png"
	path_engine_off  = "assets/warning_engine_off.png"
	path_battery_on  = "assets/warning_battery_on.png"
	path_battery_off = "assets/warning_battery_off.png"
)

type WarningLight struct {
	img_off              *ebiten.Image
	img_on               *ebiten.Image
	draw_options         *ebiten.DrawImageOptions
	is_on                bool
	temperature_key      string
	temperature_treshold float64
}

func NewWarningLight(
	path_off string,
	path_on string,
	temperature_key string,
	temperature_treshold float64,
	x float64,
	y float64,
) *WarningLight {
	img_off, _, err := ebitenutil.NewImageFromFile(path_off)
	if err != nil {
		panic(err)
	}
	img_on, _, err := ebitenutil.NewImageFromFile(path_on)
	if err != nil {
		panic(err)
	}

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(0.15, 0.15)
	opts.GeoM.Translate(x, y)
	opts.Filter = ebiten.FilterLinear
	return &WarningLight{
		img_off:              img_off,
		img_on:               img_on,
		draw_options:         opts,
		is_on:                false,
		temperature_key:      temperature_key,
		temperature_treshold: temperature_treshold,
	}
}

func NewEngineWarningLight() *WarningLight {
	pos_x := float64((window_size_x-offset_dash_x)/2 - 64)
	pos_y := float64(window_size_y - 96)
	return NewWarningLight(path_engine_off, path_engine_on, "engine_temperature_c", 80, pos_x, pos_y)
}

func NewBatteryWarningLight() *WarningLight {
	pos_x := (window_size_x-offset_dash_x)/2 + 64 - (256 * 0.2 / 2) // Correct for image width
	pos_y := float64(window_size_y - 96)
	return NewWarningLight(path_battery_off, path_battery_on, "battery_temperature_c", 80, pos_x, pos_y)
}

func (i *WarningLight) Update(kvs storage.StorageBackend) {
	isOn, err := kvs.ReadFloat64(i.temperature_key)
	if err != nil {
		panic(err)
	}
	i.is_on = isOn > i.temperature_treshold
}

func (i *WarningLight) Draw(screen *ebiten.Image) {
	if i.is_on {
		screen.DrawImage(i.img_on, i.draw_options)
	} else {
		screen.DrawImage(i.img_off, i.draw_options)
	}
}

// Not needed for now
func (i *WarningLight) IsWithinBounds(x, y int) bool { return false }
func (i *WarningLight) HandleLeftClick(x, y int)     {}
func (i *WarningLight) HandleLeftDown(x, y int)      {}
func (i *WarningLight) HandleLeftUp(x, y int)        {}
func (i *WarningLight) HandleMouseEnter(x, y int)    {}
func (i *WarningLight) HandleMouseLeave(x, y int)    {}
