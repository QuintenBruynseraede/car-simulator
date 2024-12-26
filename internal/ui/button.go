package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/user/car-simulator/internal/storage"
)

type Button struct {
	text string

	posX   float64
	posY   float64
	width  float64
	height float64

	color            color.Color
	colorHover       color.Color
	textColor        color.Color
	textColorHover   color.Color
	currentColor     color.Color
	currentTextColor color.Color

	action func()
}

func NewButton(text string, pos_x, pos_y, width, height float64, action_fn func()) *Button {
	return &Button{
		text:             text,
		posX:             pos_x,
		posY:             pos_y,
		width:            width,
		height:           height,
		color:            light_gray,
		colorHover:       gray,
		textColor:        white,
		textColorHover:   black,
		currentColor:     light_gray,
		currentTextColor: white,
		action:           action_fn,
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(b.posX), float32(b.posY), float32(b.width), float32(b.height), b.currentColor, false)
	ebitenutil.DebugPrintAt(screen, b.text, int(b.posX+10), int(b.posY+8))
}

func (b *Button) IsWithinBounds(x, y int) bool {
	return x >= int(b.posX) && x <= int(b.posX+b.width) && y >= int(b.posY) && y <= int(b.posY+b.height)
}

func (b *Button) HandleLeftClick(_, _ int) {
	b.action()
}

func (b *Button) HandleLeftDown(_, _ int) {}
func (b *Button) HandleLeftUp(_, _ int)   {}

func (b *Button) Update(kvs storage.StorageBackend) {}

func (b *Button) HandleMouseEnter(_, _ int) {
	b.currentColor = b.colorHover
	b.currentTextColor = b.textColorHover
}

func (b *Button) HandleMouseLeave(_, _ int) {
	b.currentColor = b.color
	b.currentTextColor = b.textColor
}
