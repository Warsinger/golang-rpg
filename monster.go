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

func (m *Monster) Draw(screen *ebiten.Image) {
	if m.Alive() {
		offset := m.Size / 2
		vector.DrawFilledRect(screen, m.X-offset, m.Y-offset, m.Size, m.Size, color.RGBA{255, 0, 0, 255}, true)
	}
	m.DrawInfo(screen, m.TextOffset)

}

func (m *Monster) Select(screen *ebiten.Image) {
	if m.Alive() {
		offset := m.Size / 2
		vector.StrokeRect(screen, m.X-offset, m.Y-offset, m.Size, m.Size, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}

func (e *Monster) TextOffset() (float32, float32) {
	x := e.X - e.Size/2
	y := e.Y - e.Size/2
	return x, y
}

func (m *Monster) Loot() Usable {
	return &Treasure{Item: Item{Value: m.Gold, Object: Object{m.X + m.Size/2, m.Y + m.Size/2, 15}}}
}
