package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Item struct {
	Object
	Value int
	Used  bool
}

type Usable interface {
	Draw(screen *ebiten.Image, b *Board)
	Select(screen *ebiten.Image, b *Board)

	Use(e *Entity)
	inRange(e *Entity, reach int) bool
	Refresh()
	GetObject() *Object
}

type Treasure struct {
	Item
}

type HealthPack struct {
	Item
}

func (h *HealthPack) Use(e *Entity) {
	if !h.Used {
		e.CurrentHealth = int(math.Min(float64(e.CurrentHealth+h.Value), float64(e.MaxHealth)))
		h.Used = true
	}
}

func (t *Treasure) Use(e *Entity) {
	if !t.Used {
		e.Gold += t.Value
		t.Used = true
	}
}

func (i *Item) Refresh() {
	i.Used = false
}

func (i *Item) GetObject() *Object {
	return &i.Object
}

func (i *Item) inRange(e *Entity, reach int) bool {
	return inRange(&i.Object, &e.Object, reach)
}

func (t *Treasure) Draw(screen *ebiten.Image, b *Board) {
	drawItem(screen, &t.Item, color.RGBA{255, 215, 0, 255}, b)

	// t.DrawInfo(screen, t.TextOffset)

}
func (h *HealthPack) Draw(screen *ebiten.Image, b *Board) {
	drawItem(screen, &h.Item, color.RGBA{100, 255, 100, 255}, b)

	// t.DrawInfo(screen, t.TextOffset)
}

func drawItem(screen *ebiten.Image, i *Item, c color.Color, b *Board) {
	if !i.Used {
		x, y := b.GridToXY(i.GridX, i.GridY)
		s := float32(i.Size * b.GridSize)
		vector.DrawFilledRect(screen, x, y, s, s, c, true)
	}
}

func (i *Item) Select(screen *ebiten.Image, b *Board) {
	if !i.Used {
		x, y := b.GridToXY(i.GridX, i.GridY)
		s := float32(i.Size * b.GridSize)
		vector.StrokeRect(screen, x, y, s, s, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}
