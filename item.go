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
	Draw(screen *ebiten.Image)
	Select(screen *ebiten.Image)

	Use(e *Entity)
	inRange(e *Entity) bool
	Refresh()
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

func (i *Item) inRange(e *Entity) bool {
	return inRange(&i.Object, &e.Object)
}

func (t *Treasure) Draw(screen *ebiten.Image) {
	drawItem(screen, &t.Item, color.RGBA{255, 215, 0, 255})

	// t.DrawInfo(screen, t.TextOffset)

}
func (h *HealthPack) Draw(screen *ebiten.Image) {
	drawItem(screen, &h.Item, color.RGBA{100, 255, 100, 255})

	// t.DrawInfo(screen, t.TextOffset)
}

func drawItem(screen *ebiten.Image, i *Item, c color.Color) {
	if !i.Used {
		offset := i.Size / 2
		vector.DrawFilledRect(screen, i.X-offset, i.Y-offset, i.Size, i.Size, c, true)
	}
}

func (i *Item) Select(screen *ebiten.Image) {
	if !i.Used {
		offset := i.Size / 2
		vector.StrokeRect(screen, i.X-offset, i.Y-offset, i.Size, i.Size, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}
