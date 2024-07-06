package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Monster struct {
	Entity
	Attacker
}

func (m *Monster) Draw(screen *ebiten.Image, b *Board) {
	x, y := gridToXY(m.GridX, m.GridY, b)
	if m.Alive() {
		s := float32(m.Size * b.GridSize)
		vector.DrawFilledRect(screen, x, y, s, s, color.RGBA{255, 0, 0, 255}, true)
	}
	m.DrawInfo(screen, x, y)

}

func (m *Monster) Select(screen *ebiten.Image, b *Board) {
	if m.Alive() {
		x, y := gridToXY(m.GridX, m.GridY, b)
		s := float32(m.Size * b.GridSize)
		vector.StrokeRect(screen, x, y, s, s, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}

func (m *Monster) Loot() Usable {
	return &Treasure{Item: Item{Value: m.Gold, Object: Object{m.GridX + 1, m.GridY, 1}}}
}
