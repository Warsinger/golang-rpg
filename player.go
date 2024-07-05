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
	vector.DrawFilledCircle(screen, p.X, p.Y, p.Size, c, true)
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

func (p *Player) DrawInfo(screen *ebiten.Image, textOffset func() (float32, float32)) {
	// Draw health inside the character
	var infoText string
	if p.CurrentHealth > 0 {
		infoText = fmt.Sprintf("%s(%d)\n%d/%d\n%dg\n", p.Name, p.Level, p.CurrentHealth, p.MaxHealth, p.Gold)
	} else {
		infoText = fmt.Sprintf("Dead\n%s", p.Name)
	}
	x, y := textOffset()
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = mplusNormalFace.Size * 1.5
	text.Draw(screen, infoText, mplusNormalFace, op)
}
