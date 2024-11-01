package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/user/car-simulator/internal/storage"
)

const (
	image_path_left_on   = "assets/indicator_left_on.png"
	image_path_left_off  = "assets/indicator_left_off.png"
	image_path_right_on  = "assets/indicator_right_on.png"
	image_path_right_off = "assets/indicator_right_off.png"
)

type Indicator struct {
	img_off      *ebiten.Image
	img_on       *ebiten.Image
	draw_options *ebiten.DrawImageOptions
	side         string
	is_on        bool
}

func NewIndicatorLeft() *Indicator {
	img_off, _, err := ebitenutil.NewImageFromFile(image_path_left_off)
	if err != nil {
		panic(err)
	}
	img_on, _, err := ebitenutil.NewImageFromFile(image_path_left_on)
	if err != nil {
		panic(err)
	}
	x := (window_size_x-offset_dash_x)/2 - 256
	y := offset_dash_y + 32

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(0.15, 0.15)
	opts.GeoM.Translate(float64(x), float64(y))
	opts.Filter = ebiten.FilterLinear
	return &Indicator{
		img_off:      img_off,
		img_on:       img_on,
		draw_options: opts,
		side:         "left",
	}
}

func NewIndicatorRight() *Indicator {
	img_off, _, err := ebitenutil.NewImageFromFile(image_path_right_off)
	if err != nil {
		panic(err)
	}
	img_on, _, err := ebitenutil.NewImageFromFile(image_path_right_on)
	if err != nil {
		panic(err)
	}

	x := (window_size_x-offset_dash_x)/2 + 256
	y := offset_dash_y + 32
	scale := 0.15

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(0.15, 0.15)
	opts.GeoM.Translate(float64(x), float64(y))
	opts.GeoM.Translate(float64(-img_on.Bounds().Dx())*scale, 0) // Correct for image width
	opts.Filter = ebiten.FilterLinear
	return &Indicator{
		img_off:      img_off,
		img_on:       img_on,
		draw_options: opts,
		side:         "right",
	}
}

func (i *Indicator) Update(kvs storage.StorageBackend) {
	if i.side == "left" {
		i.is_on = kvs.ReadString("indicator_left_status") == "on"
	} else if i.side == "right" {
		i.is_on = kvs.ReadString("indicator_right_status") == "on"
	}
}

func (i *Indicator) Draw(screen *ebiten.Image) {
	if i.is_on {
		screen.DrawImage(i.img_on, i.draw_options)
	} else {
		screen.DrawImage(i.img_off, i.draw_options)
	}
}

// Not needed for now
func (g *Indicator) IsWithinBounds(x, y int) bool { return false }
func (g *Indicator) HandleLeftClick(x, y int)     {}
func (g *Indicator) HandleMouseEnter(x, y int)    {}
func (g *Indicator) HandleMouseLeave(x, y int)    {}
