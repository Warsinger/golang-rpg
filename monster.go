package main

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gopkg.in/yaml.v3"
)

type MonsterInfo struct {
	EntityInfo
	AttackInfo
}
type Monster interface {
	Entity
	Attacker

	Draw(screen *ebiten.Image, b Board)
	Select(screen *ebiten.Image, b Board)
	Loot() Item
}

func (m *MonsterInfo) Draw(screen *ebiten.Image, b Board) {
	x, y := b.GridToXY(m.GridX, m.GridY)
	if m.Alive() {
		s := float32(m.Size * b.GetGridSize())
		vector.DrawFilledRect(screen, x, y, s, s, color.RGBA{255, 0, 0, 255}, true)
	}
	m.DrawInfo(screen, x, y)

}

func (m *MonsterInfo) Select(screen *ebiten.Image, b Board) {
	if m.Alive() {
		x, y := b.GridToXY(m.GridX, m.GridY)
		s := float32(m.Size * b.GetGridSize())
		vector.StrokeRect(screen, x, y, s, s, 2, color.RGBA{0, 255, 255, 255}, true)
	}
}

func (m *MonsterInfo) Loot() Item {
	return &TreasureInfo{ItemInfo: ItemInfo{Value: m.Gold, ObjectInfo: ObjectInfo{m.GridX + 1, m.GridY, 1}}}
}

func LoadMonsters(b Board) ([]MonsterInfo, error) {
	yamlFile, err := os.ReadFile("config/monsters.yml")
	if err != nil {
		return nil, err
	}
	var monsters []MonsterInfo
	err = yaml.Unmarshal(yamlFile, &monsters)
	if err != nil {
		return nil, err
	}

	// place monsters on the board
	for _, m := range monsters {
		b.AddObjectToBoard(&m)
	}

	return monsters, nil
}
