package ui

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/user/car-simulator/internal/storage"
)

const (
	font_file = "assets/Open_24_Display_St.ttf"
)

type EightSegmentDisplay struct {
	value int
	x     float64
	y     float64
	font  *text.GoTextFaceSource
	color color.RGBA
}

func NewEightSegmentDisplay(x float64, y float64) *EightSegmentDisplay {
	fontBytes, err := os.ReadFile(font_file)
	font, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatal(err)
	}

	return &EightSegmentDisplay{
		value: 0,
		x:     x,
		y:     y,
		font:  font,
		color: color.RGBA{190, 190, 190, 255},
	}
}

func (i *EightSegmentDisplay) Update(kvs storage.StorageBackend) {
	value, err := kvs.ReadInt("gear")
	if err != nil {
		panic(err)
	}
	i.value = value
}

func (i *EightSegmentDisplay) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(i.x, i.y)
	op.ColorScale.ScaleWithColor(i.color)
	text.Draw(screen, fmt.Sprintf("%v", i.value), &text.GoTextFace{
		Source: i.font,
		Size:   40,
	}, op)
}

// Not needed for now
func (i *EightSegmentDisplay) IsWithinBounds(x, y int) bool { return false }
func (i *EightSegmentDisplay) HandleLeftClick(x, y int)     {}
func (i *EightSegmentDisplay) HandleMouseEnter(x, y int)    {}
func (i *EightSegmentDisplay) HandleMouseLeave(x, y int)    {}
func (i *EightSegmentDisplay) HandleLeftDown(_, _ int)      {}
func (i *EightSegmentDisplay) HandleLeftUp(_, _ int)        {}
