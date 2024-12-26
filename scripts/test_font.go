package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand/v2"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var mplusFaceSource *text.GoTextFaceSource

func init() {
	fontBytes, err := os.ReadFile("assets/Open_24_Display_St.ttf")
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

type Game struct {
	counter   int
	text      string
	textColor color.RGBA
}

func (g *Game) Update() error {
	// Change the text color for each second.
	if g.counter%ebiten.TPS() == 0 {
		g.text = "123456^"

		g.textColor.R = 0x80 + uint8(rand.IntN(0x7f))
		g.textColor.G = 0x80 + uint8(rand.IntN(0x7f))
		g.textColor.B = 0x80 + uint8(rand.IntN(0x7f))
		g.textColor.A = 0xff
	}
	g.counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	const (
		normalFontSize = 24
		bigFontSize    = 48
	)

	const x = 20

	// Draw
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, 110)
	op.ColorScale.ScaleWithColor(g.textColor)
	op.LineSpacing = bigFontSize * 1.2
	text.Draw(screen, g.text, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   bigFontSize,
	}, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Font (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
