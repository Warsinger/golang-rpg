package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Item struct {
	Object
	value int
	used  bool
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

type Health struct {
	Item
}

func (h *Health) Use(e *Entity) {
	if !h.used {
		e.currentHealth = int(math.Min(float64(e.currentHealth+h.value), float64(e.maxHealth)))
		h.used = true
	}
}

func (t *Treasure) Use(e *Entity) {
	if !t.used {
		e.gold += t.value
		t.used = true
	}
}

func (i *Item) Refresh() {
	i.used = false
}

func (i *Item) inRange(e *Entity) bool {
	return inRange(&i.Object, &e.Object)
}

func (t *Treasure) Draw(screen *ebiten.Image) {
	drawItem(screen, &t.Item, color.RGBA{255, 215, 0, 255})

	// t.DrawInfo(screen, t.TextOffset)

}
func (h *Health) Draw(screen *ebiten.Image) {
	drawItem(screen, &h.Item, color.RGBA{100, 255, 100, 255})

	// t.DrawInfo(screen, t.TextOffset)
}

func drawItem(screen *ebiten.Image, i *Item, c color.Color) {
	if !i.used {
		offset := i.size / 2
		vector.DrawFilledRect(screen, i.x-offset, i.y-offset, i.size, i.size, c, true)
	}
}

func (i *Item) Select(screen *ebiten.Image) {
	if !i.used {
		offset := i.size / 2
		vector.StrokeRect(screen, i.x-offset, i.y-offset, i.size, i.size, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}
