package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/user/car-simulator/internal/storage"
)

type UIObject interface {
	Update(kvs storage.StorageBackend)
	Draw(screen *ebiten.Image)
	IsWithinBounds(x, y int) bool
	HandleLeftClick(x, y int)
	HandleLeftDown(x, y int)
	HandleLeftUp(x, y int)
	HandleMouseEnter(x, y int)
	HandleMouseLeave(x, y int)
}
