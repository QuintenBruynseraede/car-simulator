package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	img          *ebiten.Image
	draw_options *ebiten.DrawImageOptions
)

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("assets/indicator_left_on.png")
	draw_options = &ebiten.DrawImageOptions{}
	draw_options.GeoM.Scale(0.1, 0.1)
	draw_options.GeoM.Translate(32, 32)
	draw_options.Filter = ebiten.FilterLinear
	if err != nil {
		log.Fatal(err)
	}
}

type UI struct{}

func (g *UI) Update() error {
	return nil
}

func (g *UI) Draw(screen *ebiten.Image) {
	screen.DrawImage(img, draw_options)
}

func (g *UI) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Car Simulator")
	ebiten.SetVsyncEnabled(true)
	if err := ebiten.RunGame(&UI{}); err != nil {
		log.Fatal(err)
	}
}
