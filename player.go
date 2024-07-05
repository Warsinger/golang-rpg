package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	Entity
	Attacker
}

func (p *Player) TextOffset() (float32, float32) {
	// x := p.x - p.size*2
	// y := p.y - p.size/2
	// return x, y
	return 4, 4
}

func (p *Player) Draw(screen *ebiten.Image) {
	var c color.Color
	if p.Alive() {
		c = color.RGBA{0, 255, 0, 255}
	} else {
		c = color.RGBA{128, 128, 128, 255}
	}
	vector.DrawFilledCircle(screen, p.x, p.y, p.size, c, true)
	p.DrawInfo(screen, p.TextOffset)
}
