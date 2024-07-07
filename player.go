package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gopkg.in/yaml.v3"
)

type PlayerInfo struct {
	EntityInfo
	AttackInfo
	Experience int
	Level      int
}

type Player interface {
	Entity
	Attacker

	Draw(screen *ebiten.Image, b Board)

	GetExperience() int
	GetLevel() int
	AttackMonster(m Monster)
	AddExperience(xp int)
	LevelUp(newLevel int)
	Move(direction Direction, b Board) bool
}

func LoadPlayer(b Board) (Player, error) {
	yamlFile, err := os.ReadFile("config/player.yml")
	if err != nil {
		return nil, err
	}
	var player PlayerInfo
	err = yaml.Unmarshal(yamlFile, &player)
	if err != nil {
		return nil, err
	}

	b.AddObjectToBoard(&player)
	return &player, nil
}

func (p *PlayerInfo) Draw(screen *ebiten.Image, b Board) {
	var c color.Color
	if p.Alive() {
		c = color.RGBA{0, 255, 0, 255}
	} else {
		c = color.RGBA{128, 128, 128, 255}
	}
	x, y := b.GridToXY(p.GridX, p.GridY)
	x += float32(b.GetGridSize()*p.Size) / 2
	y += float32(b.GetGridSize()*p.Size) / 2
	r := float32(p.Size*b.GetGridSize()) / 2
	vector.DrawFilledCircle(screen, x, y, r, c, true)
	p.DrawInfo(screen, 4, 4)
}

func (p *PlayerInfo) AttackMonster(m Monster) {
	p.Attack(m)
	if m.Alive() {
		// if monster is still alive calculate the monster's attack value and subtract from player's health
		m.Attack(p)
	}
}

func (p *PlayerInfo) AddExperience(xp int) {
	p.Experience += xp
	newLevel := p.Experience / 25
	if newLevel > p.Level {
		p.LevelUp(newLevel)
	}
}

func (p *PlayerInfo) LevelUp(newLevel int) {
	p.AttackPower++
	p.Level = newLevel
	fmt.Println("Level Up!")
}

func (p *PlayerInfo) DrawInfo(screen *ebiten.Image, x, y float32) {
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

func (p *PlayerInfo) Move(direction Direction, b Board) bool {
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

	if b.CanOccupySpace(p, gx, gy) {
		p.GridX = gx
		p.GridY = gy
		return true
	}
	return false
}

func (p *PlayerInfo) GetExperience() int {
	return p.Experience
}
func (p *PlayerInfo) GetLevel() int {
	return p.Level
}
