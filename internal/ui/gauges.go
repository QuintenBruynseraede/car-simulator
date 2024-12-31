package ui

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/user/car-simulator/internal/storage"
)

type Gauge struct {
	value_key string

	// Gauge
	background_image *ebiten.Image
	scale            float64
	pos_x            float64
	pos_y            float64

	// Needle
	needle_image *ebiten.Image
	needle_scale float64
	needle_pos_x float64
	needle_pos_y float64

	// For gauge needle
	angle           float64
	min_value_angle float64
	max_value_angle float64
	min_value       float64
	max_value       float64
}

func NewSpedometer() *Gauge {
	img, _, _ := ebitenutil.NewImageFromFile("assets/spedometer.png")    // TODO handle error =3rd return value
	needle_img, _, _ := ebitenutil.NewImageFromFile("assets/needle.png") // TODO handle error =3rd return value
	value_key := "velocity_kmh"

	return &Gauge{
		value_key:        value_key,
		background_image: img,
		scale:            0.12,
		pos_x:            100,
		pos_y:            550,
		needle_image:     needle_img,
		needle_scale:     0.08,
		needle_pos_x:     150 + 57,
		needle_pos_y:     695 - 39,
		angle:            0,
		min_value_angle:  4*math.Pi/5 + 0.12, // Magic number that points the needle to the zero value
		max_value_angle:  1 * math.Pi / 5,
		min_value:        0.0,
		max_value:        270.0,
	}
}

func NewRPMGauge() *Gauge {
	img, _, _ := ebitenutil.NewImageFromFile("assets/rpm_gauge.png")     // TODO handle error =3rd return value
	needle_img, _, _ := ebitenutil.NewImageFromFile("assets/needle.png") // TODO handle error =3rd return value
	value_key := "engine_rpm"

	return &Gauge{
		value_key:        value_key,
		background_image: img,
		scale:            0.12,
		pos_x:            700,
		pos_y:            550,
		needle_image:     needle_img,
		needle_scale:     0.08,
		needle_pos_x:     600 + 150 + 57,
		needle_pos_y:     695 - 39,
		angle:            0,
		min_value_angle:  4*math.Pi/5 + 0.12, // Magic number that points the needle to the zero value
		max_value_angle:  1 * math.Pi / 5,
		min_value:        0.0,
		max_value:        8000,
	}
}

func (g *Gauge) Update(kvs storage.StorageBackend) {
	speed, err := kvs.ReadFloat64(g.value_key)
	if err != nil {
		panic(err)
	}

	perc := (speed - g.min_value) / (g.max_value - g.min_value)
	g.angle = g.min_value_angle + perc*(2*math.Pi+g.max_value_angle-g.min_value_angle)
}

func (g *Gauge) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(g.scale, g.scale)
	opts.GeoM.Translate(g.pos_x, g.pos_y)
	opts.Filter = ebiten.FilterLinear

	needle_opts := &ebiten.DrawImageOptions{}
	offset_x := float64(g.needle_image.Bounds().Size().X) / 2.0
	offset_y := float64(g.needle_image.Bounds().Size().Y) / 2.0
	needle_opts.GeoM.Translate(-offset_x, -offset_y) // Center before rotation

	needle_opts.GeoM.Rotate(g.angle)
	needle_opts.GeoM.Scale(g.needle_scale, g.needle_scale)
	needle_opts.GeoM.Translate(g.needle_pos_x, g.needle_pos_y)
	needle_opts.Filter = ebiten.FilterLinear

	screen.DrawImage(g.background_image, opts)
	screen.DrawImage(g.needle_image, needle_opts)
}

// Not needed for now
func (g *Gauge) IsWithinBounds(x, y int) bool { return false }
func (g *Gauge) HandleLeftClick(x, y int)     {}
func (g *Gauge) HandleMouseEnter(x, y int)    {}
func (g *Gauge) HandleMouseLeave(x, y int)    {}
func (g *Gauge) HandleLeftDown(_, _ int)      {}
func (g *Gauge) HandleLeftUp(_, _ int)        {}
