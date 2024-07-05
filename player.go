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
	experience int
	level      int
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

func (p *Player) AttackMonster(m *Monster) {
	p.Attack(&m.Entity)
	if m.Alive() {
		// if monster is still alive calculate the monster's attack value and subtract from player's health
		m.Attack(&p.Entity)
	}
}

func (p *Player) AddExperience(xp int) {
	p.experience += xp
	newLevel := p.experience / 25
	if newLevel > p.level {
		p.LevelUp(newLevel)
	}
}

func (p *Player) LevelUp(newLevel int) {
	p.attack++
	p.level = newLevel
	fmt.Println("Level Up!")
}

func (p *Player) DrawInfo(screen *ebiten.Image, textOffset func() (float32, float32)) {
	// Draw health inside the character
	var infoText string
	if p.currentHealth > 0 {
		infoText = fmt.Sprintf("%s(%d)\n%d/%d\n%dg\n", p.name, p.level, p.currentHealth, p.maxHealth, p.gold)
	} else {
		infoText = fmt.Sprintf("Dead\n%s", p.name)
	}
	x, y := textOffset()
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = mplusNormalFace.Size * 1.5
	text.Draw(screen, infoText, mplusNormalFace, op)
}
