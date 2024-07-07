package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	Entity
	Attacker
	Experience int
	Level      int
}

func (p *Player) Draw(screen *ebiten.Image, b *Board) {
	var c color.Color
	if p.Alive() {
		c = color.RGBA{0, 255, 0, 255}
	} else {
		c = color.RGBA{128, 128, 128, 255}
	}
	x, y := b.GridToXY(p.GridX, p.GridY)
	x += float32(b.GridSize*p.Size) / 2
	y += float32(b.GridSize*p.Size) / 2
	r := float32(p.Size*b.GridSize) / 2
	vector.DrawFilledCircle(screen, x, y, r, c, true)
	p.DrawInfo(screen, 4, 4)
}

func (p *Player) AttackMonster(m *Monster) {
	p.Attack(&m.Entity)
	if m.Alive() {
		// if monster is still alive calculate the monster's attack value and subtract from player's health
		m.Attack(&p.Entity)
	}
}

func (p *Player) AddExperience(xp int) {
	p.Experience += xp
	newLevel := p.Experience / 25
	if newLevel > p.Level {
		p.LevelUp(newLevel)
	}
}

func (p *Player) LevelUp(newLevel int) {
	p.AttackPower++
	p.Level = newLevel
	fmt.Println("Level Up!")
}

func (p *Player) DrawInfo(screen *ebiten.Image, x, y float32) {
	var infoText string
	if p.CurrentHealth > 0 {
		infoText = fmt.Sprintf("%s(%d)\n%d/%d\n%dg\n", p.Name, p.Level, p.CurrentHealth, p.MaxHealth, p.Gold)
	} else {
		infoText = fmt.Sprintf("Dead\n%s", p.Name)
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = mplusNormalFace.Size * 1.5
	text.Draw(screen, infoText, mplusNormalFace, op)
}

func (p *Player) Move(direction Direction, b *Board) bool {
	gx, gy := p.GridX, p.GridY
	switch direction {
	case Up:
		gy -= 1
	case Down:
		gy += 1
	case Right:
		gx += 1
	case Left:
		gx -= 1
	}

	if b.CanOccupySpace(&p.Object, gx, gy) {
		p.GridX = gx
		p.GridY = gy
		return true
	}
	return false
}
