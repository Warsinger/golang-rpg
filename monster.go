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
		offset := m.size / 2
		vector.DrawFilledRect(screen, m.x-offset, m.y-offset, m.size, m.size, color.RGBA{255, 0, 0, 255}, true)

	}
	m.DrawInfo(screen, m.TextOffset)

}

func (m *Monster) Select(screen *ebiten.Image) {
	if m.Alive() {
		offset := m.size / 2
		vector.StrokeRect(screen, m.x-offset, m.y-offset, m.size, m.size, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}

func (e *Monster) TextOffset() (float32, float32) {
	x := e.x - e.size/2
	y := e.y - e.size/2
	return x, y
}
