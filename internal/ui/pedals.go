package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/user/car-simulator/internal/storage"
)

type Pedal struct {
	posX  float64
	posY  float64
	image *ebiten.Image
	scale float64

	action_down func()
	action_up   func()

	// Dynamic
	is_down bool
}

func NewGasPedal(action_down func(), action_up func()) *Pedal {
	img, _, err := ebitenutil.NewImageFromFile("assets/gas_pedal.png")
	if err != nil {
		panic(err)
	}
	return &Pedal{
		posX:        770,
		posY:        240,
		image:       img,
		scale:       0.1,
		action_down: action_down,
		action_up:   action_up,

		is_down: false,
	}
}

func NewBrakePedal(action_down func(), action_up func()) *Pedal {
	img, _, err := ebitenutil.NewImageFromFile("assets/brake_pedal.png")
	if err != nil {
		panic(err)
	}
	return &Pedal{
		posX:        670,
		posY:        240,
		image:       img,
		scale:       0.1,
		action_down: action_down,
		action_up:   action_up,

		is_down: false,
	}
}

func (p *Pedal) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	if p.is_down {
		opts.ColorM.Scale(0.5, 0.5, 0.5, 1)
		opts.GeoM.Scale(1, 0.85)
	}
	opts.GeoM.Scale(p.scale, p.scale)
	opts.GeoM.Translate(p.posX, p.posY)
	opts.Filter = ebiten.FilterLinear

	screen.DrawImage(p.image, opts)
}

func (p *Pedal) IsWithinBounds(x, y int) bool {
	w, h := float64(p.image.Bounds().Bounds().Size().X), float64(p.image.Bounds().Bounds().Size().Y)

	return (float64(x) >= p.posX &&
		float64(x) <= p.posX+w*p.scale &&
		float64(y) >= p.posY &&
		float64(y) <= p.posY+h*p.scale)
}

func (p *Pedal) HandleLeftClick(x, y int) {
	p.is_down = true
	p.action_down()
}

func (p *Pedal) HandleLeftDown(x, y int) {}
func (p *Pedal) HandleLeftUp(x, y int) {
	p.is_down = false
	p.action_up()
}

func (p *Pedal) Update(kvs storage.StorageBackend) {}

func (p *Pedal) HandleMouseEnter(x, y int) {}

func (p *Pedal) HandleMouseLeave(x, y int) {
	p.is_down = false
	p.action_up()
}
